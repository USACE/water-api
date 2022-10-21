package models

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/USACE/water-api/api/app"
	"github.com/USACE/water-api/api/interpolation"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UsgsRating struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Name   string    `json:"name" db:"name"`
	Method string    `json:"method" db:"method"`
	S3Key  string    `json:"s3key" db:"s3key"`
}

func IndepDepTable(db *pgxpool.Pool, n string, m string) ([]float64, []float64, error) {
	var r UsgsRating
	if err := pgxscan.Get(
		context.Background(),
		db,
		&r,
		`SELECT * FROM usgs_rating WHERE name = $1 AND "method" = $2`,
		n, m,
	); err != nil {
		return nil, nil, err
	}

	cfg, err := app.GetConfig()
	if err != nil {
		return nil, nil, err
	}
	s3Config := app.AWSConfig(*cfg)

	sess, err := session.NewSession(&s3Config)
	if err != nil {
		return nil, nil, err
	}
	svc := s3.New(sess)

	params := &s3.SelectObjectContentInput{
		Bucket:         aws.String(cfg.AWSS3Bucket),
		Key:            aws.String(r.S3Key),
		ExpressionType: aws.String(s3.ExpressionTypeSql),
		Expression:     aws.String(fmt.Sprintf("SELECT * FROM S3Object")),
		InputSerialization: &s3.InputSerialization{
			CSV: &s3.CSVInput{
				FileHeaderInfo: aws.String(s3.FileHeaderInfoUse),
			},
			// CompressionType: aws.String(s3.CompressionTypeGzip),
		},
		OutputSerialization: &s3.OutputSerialization{
			CSV: &s3.CSVOutput{},
		},
	}

	resp, err := svc.SelectObjectContent(params)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, fmt.Errorf("failed to read from SelectObjectContent EventStream, %v", err)
	}

	resReader := csv.NewReader(results)

	xa := []float64{}
	ya := []float64{}
	for {
		r, err := resReader.Read()
		if len(r) > 0 {
			indep, _ := strconv.ParseFloat(r[0], 32)
			dep, _ := strconv.ParseFloat(r[1], 32)
			xa = append(xa, indep)
			ya = append(ya, dep)
		}
		if err == io.EOF {
			break
		}
	}
	return xa, ya, nil
}

// Rate1Var
func Rate1Var(db *pgxpool.Pool, n string, m string, x float64) (float64, error) {
	var (
		err    error
		xa, ya []float64
		y      float64
	)
	if xa, ya, err = IndepDepTable(db, n, m); err != nil {
		return 0, err
	}
	if y, err = interpolation.LinearEval(xa, ya, x, &interpolation.InterpAccel{}); err != nil {
		return 0, err
	}

	// return y, nil
	return y, nil
}
