package migration

import "database/sql"

func CreateCounterTable(db *sql.DB) error {
	query := `
create table counter_metric
(
    id        integer not null
        constraint counter_metric_pk
            primary key,
    type      varchar(255),
    value     int8
);

`

	_, err := db.Exec(query)

	return err
}

func CreateGaugeTable(db *sql.DB) error {
	query := `
create table gauge_metric
(
	id        integer not null
        constraint gauge_metric_pk
            primary key,
    type      varchar(255),
    value     double precision
);
`

	_, err := db.Exec(query)

	return err
}
