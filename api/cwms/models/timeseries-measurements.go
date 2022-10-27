package models

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/USACE/water-api/api/app"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Measurements struct {
	Times  []time.Time `json:"times,omitempty"`
	Values []float64   `json:"values,omitempty"`
}

func (m Measurements) LatestTime() *time.Time {
	if len(m.Times) > 0 {
		return &m.Times[len(m.Times)-1]
	}
	return nil
}

func (m Measurements) LatestValue() *float64 {
	if len(m.Values) > 0 {
		return &m.Values[len(m.Values)-1]
	}
	return nil
}

func CreateOrUpdateTimeseriesMeasurements(db *pgxpool.Pool, c TimeseriesCollection) ([]Timeseries, error) {

	tx, err := db.Begin(context.Background())
	if err != nil {
		return make([]Timeseries, 0), err
	}
	defer tx.Rollback(context.Background())

	updatedIDs := make([]uuid.UUID, 0)

	for _, t := range c.Items {

		// If measurements aren't properly provided, skip this item
		if t.Measurements == nil || t.Measurements.Times == nil || t.Measurements.Values == nil {
			// fmt.Println("skipping " + t.Key)
			continue
		}

		rows, err := tx.Query(
			context.Background(),
			`UPDATE timeseries 
			SET latest_time = $4, latest_value = $5
			WHERE datasource_key = $3
			AND datasource_id = (
				SELECT d.id FROM datasource d 
				JOIN datasource_type dt ON dt.id = d.datasource_type_id 
				JOIN provider p ON p.id = d.provider_id 
				WHERE lower(p.slug) = lower($2) AND lower(dt.slug) = lower($1)
			)
			RETURNING id`,
			t.DatasourceType, t.Provider, t.Key, t.Measurements.LatestTime(), t.Measurements.LatestValue(),
		)
		if err != nil {
			return make([]Timeseries, 0), err
		}
		var id uuid.UUID
		if err := pgxscan.ScanOne(&id, rows); err != nil {
			tx.Rollback(context.Background())
			return c.Items, err
		} else {
			updatedIDs = append(updatedIDs, id)
		}

	}
	tx.Commit(context.Background())

	if len(updatedIDs) == 0 {
		return nil, errors.New("no records updated")
	}

	return nil, nil
}

// GetTimeseriesMeasurements
func GetTimeseriesMeasurements(db *pgxpool.Pool, k string, a string, b string) (Measurements, error) {
	ms := Measurements{}
	fmt.Printf("%s\n%s\n%s", k, a, b)
	cfg, err := app.GetConfig()
	if err != nil {
		return ms, err
	}
	s3Config := app.AWSConfig(*cfg)

	sess, err := session.NewSession(&s3Config)
	if err != nil {
		return ms, err
	}
	svc := s3.New(sess)

	params := &s3.SelectObjectContentInput{
		Bucket:         aws.String(cfg.AWSS3Bucket),
		Key:            aws.String(k),
		ExpressionType: aws.String(s3.ExpressionTypeSql),
		Expression:     aws.String(fmt.Sprintf("SELECT * FROM S3Object s WHERE s.Date >= '%s' AND s.Date <='%s'", a, b)),
		InputSerialization: &s3.InputSerialization{
			CSV: &s3.CSVInput{
				FileHeaderInfo: aws.String(s3.FileHeaderInfoUse),
			},
			CompressionType: aws.String(s3.CompressionTypeGzip),
		},
		OutputSerialization: &s3.OutputSerialization{
			CSV: &s3.CSVOutput{},
		},
	}
	resp, err := svc.SelectObjectContent(params)
	if err != nil {
		return ms, err
	}
	defer resp.EventStream.Close()

	results, resultWriter := io.Pipe()
	go func() {
		defer resultWriter.Close()
		for event := range resp.EventStream.Events() {
			switch e := event.(type) {
			case *s3.RecordsEvent:
				resultWriter.Write(e.Payload)
				// fmt.Printf("Payload: %v\n", string(e.Payload))
			case *s3.StatsEvent:
				fmt.Printf("Processed %d bytes\n", *e.Details.BytesProcessed)
			}
		}
	}()

	if err := resp.EventStream.Err(); err != nil {
		return ms, fmt.Errorf("failed to read from SelectObjectContent EventStream, %v", err)
	}

	resReader := csv.NewReader(results)

	for {
		r, err := resReader.Read()
		if len(r) > 0 {
			t, _ := time.Parse(time.RFC3339, r[0])
			v, _ := strconv.ParseFloat(r[1], 32)
			ms.Times = append(ms.Times, t)
			ms.Values = append(ms.Values, v)
		}
		if err == io.EOF {
			break
		}
	}
	return ms, nil
}
