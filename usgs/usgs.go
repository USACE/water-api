package usgs

import (
	"errors"
	"net/http"
	"strings"

	"github.com/USACE/water-api/messages"
	"github.com/USACE/water-api/usgs/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/labstack/echo/v4"
)

func (s Store) ListSites(c echo.Context) error {

	// Get filter from query params state_id =
	var sf models.SiteFilter
	if err := c.Bind(&sf); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ss, err := models.ListSites(s.Connection, &sf)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ss)
}

func (s Store) CreateSites(c echo.Context) error {
	var sc models.SiteCollection
	if err := c.Bind(&sc); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Create Sites
	ss, err := models.CreateSites(s.Connection, sc.Items)
	if err != nil {
		// If Error was postgres error, return error message based on error code
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return c.JSON(
					http.StatusBadRequest,
					messages.NewMessage("Sites not created. Site information conflicts with an existing site"))
			}
		}
		// If not explicit error, return string of error message for debugging
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, ss)
}

func (s Store) SyncSites(c echo.Context) error {

	// Get existing sites for comparison
	var sf models.SiteFilter
	existingSites, err := models.ListSites(s.Connection, &sf)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Analyze sync payload
	var sc models.SiteCollection

	if err := c.Bind(&sc); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Sites to store only new sites
	new_sites := make([]models.Site, 0)
	update_sites := make([]models.Site, 0)

	sitemap := make(map[string]models.Site, len(existingSites))

	for _, s := range existingSites {
		sitemap[s.SiteInfo.UsgsId] = s
	}

	// Loop over payload items
	for _, site := range sc.Items {

		// If no key in map, we have a new site
		if existingSite, ok := sitemap[site.SiteInfo.UsgsId]; !ok {
			// payload site not found, adding to new_sites
			new_sites = append(new_sites, site)
		} else {

			if !site.IsEquivalent(existingSite) {
				update_sites = append(update_sites, site)
				// fmt.Println("update needed")
			}
		}
	}

	r := map[string]interface{}{}

	if len(update_sites) > 0 {
		// Update sites
		sites_updated, err := models.UpdateSites(s.Connection, update_sites)

		if err != nil {
			if pgxscan.NotFound(err) {
				return c.JSON(http.StatusNotFound, messages.DefaultMessageNotFound)
			}
			return c.JSON(http.StatusInternalServerError, err)
			// return c.String(http.StatusInternalServerError, err.Error())
		}

		r["updated"] = &sites_updated
		r["update_count"] = len(sites_updated)
	} else {
		r["updated"] = make([]models.Site, 0)
		r["update_count"] = 0
	}

	if len(new_sites) > 0 {
		// Create new Sites
		sites_created, err := models.CreateSites(s.Connection, new_sites)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		r["created"] = &sites_created
		r["create_count"] = len(sites_created)

	} else {
		r["created"] = make([]models.Site, 0)
		r["create_count"] = 0
	}

	return c.JSON(http.StatusAccepted, r)
}

func (s Store) ListParameters(c echo.Context) error {
	pp, err := models.ListParameters(s.Connection)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, pp)
}

// NOT USED - SAVE FOR WATERSHED/SITE/PARAM Enabled
// func (s Store) ListParametersEnabled(c echo.Context) error {
// 	pp, err := models.ListParametersEnabled(s.Connection)
// 	if err != nil {
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}
// 	return c.JSON(http.StatusOK, pp)
// }

func (s Store) CreateSiteParameters(c echo.Context) error {

	var sp models.SiteParameterCollection
	if err := c.Bind(&sp); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ss, err := models.CreateSiteParameters(s.Connection, sp.Items)
	if err != nil {
		if strings.Contains(string(err.Error()), "duplicate key value") {
			// return 422
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, ss)
}
