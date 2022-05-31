package gotile

import (
	"GoTile/config"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB
var geomColumnsCache map[string]GeometryColumns = make(map[string]GeometryColumns)

func init() {
	var err error = nil
	db, err = sqlx.Connect(config.Configiure.Database.Driver, config.Configiure.Database.Url)
	if err != nil {
		log.Fatalf("%s:%s connect fail.", config.Configiure.Database.Driver, config.Configiure.Database.Url)
	}
}

func GetMvtTileBinary(x, y, zoom int64, tablename, format string) ([]byte, error) {
	if !tileIsValid(x, y, zoom, format) {
		return nil, errors.New("request is not valid")
	}

	pbf, err := GetCache(x, y, zoom, tablename)
	if err == nil {
		return pbf, err
	}
	envelope := tileToEnvelope(x, y, zoom)
	sql, err := envelopeToSQL(envelope, tablename)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(sql).Scan(&pbf)

	if err != nil {
		return nil, err
	}

	MakeCache(x, y, zoom, pbf, tablename)
	return pbf, nil
}

func GetGeometryOfTable(tablename string) (GeometryColumns, error) {
	if value, ok := geomColumnsCache[tablename]; ok {
		return value, nil
	}

	var sql string = `select * from public."geometry_columns" where f_table_name = $1`
	rows, err := db.Queryx(sql, tablename)
	if err != nil {
		return GeometryColumns{}, err
	}

	rows.Next()
	var geometryColumns GeometryColumns
	err = rows.StructScan(&geometryColumns)

	if err != nil {
		return GeometryColumns{}, err
	}
	for rows.Next() {
		log.Fatal("found other srid in geometry_columns, use first.")
	}
	geomColumnsCache[tablename] = geometryColumns
	return geometryColumns, nil
}

func tileIsValid(x, y, zoom int64, format string) bool {
	if format != "pbf" && format != "mvt" {
		return false
	}

	size := int64(math.Pow(2, float64(zoom)))
	if x >= size || y >= size {
		return false
	}
	if x < 0 || y < 0 {
		return false
	}
	return true
}

func tileToEnvelope(x, y, z int64) [4]float64 {
	var tileSize float64 = math.Pow(2, float64(z))

	var lonMin float64 = (float64(x)/tileSize)*360.0 - 180.0
	var latRedMin float64 = math.Atan(math.Sinh(math.Pi * (1.0 - (2.0*float64(y))/tileSize)))
	var latMin float64 = (180 * latRedMin) / math.Pi

	var lonMax float64 = ((float64(x)+1)/tileSize)*360.0 - 180.0
	var latRedMax float64 = math.Atan(math.Sinh(math.Pi * (1 - (2*(float64(y)+1))/tileSize)))
	var latMax float64 = (180 * latRedMax) / math.Pi

	return [4]float64{lonMin, latMin, lonMax, latMax}
}

func envelopeToBoundsSQL(envelope [4]float64) string {
	return fmt.Sprintf("ST_MakeEnvelope(%f, %f, %f, %f, 4326)",
		envelope[0], envelope[1], envelope[2], envelope[3])
}

//Generate a SQL query to pull a tile worth of MVT data
//from the table of interest.
func envelopeToSQL(envelope [4]float64, tablename string) (string, error) {
	envelopeSQL := envelopeToBoundsSQL(envelope)
	geometryColumns, err := GetGeometryOfTable(tablename)
	if err != nil {
		return "", err
	}
	//Materialize the bounds
	//Select the relevant geometry and clip to MVT bounds
	//Convert to MVT format
	executeSQL := fmt.Sprintf(
		`with 
				bounds as (select st_transform(%s, %d) bound),
				mvtgeom as (select st_asmvtgeom(t.geom, bounds.bound, 4096, 256, true) geoms from %s t, bounds
							where t.geom && bounds.bound)
				select st_asmvt(mvtgeom, '%s' ,4096 ,'geoms') pbf from mvtgeom`,
		envelopeSQL, geometryColumns.Srid, tablename, tablename)

	return executeSQL, nil
}
