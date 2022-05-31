package gotile

type GeometryColumns struct {
	Srid           int64  `db:"srid"`
	TableName      string `db:"f_table_name"`
	TableSchema    string `db:"f_table_schema"`
	TableCatalog   string `db:"f_table_catalog"`
	GeometryColumn string `db:"f_geometry_column"`
	CoordDimension int64  `db:"coord_dimension"`
	Type           string `db:"type"`
}
