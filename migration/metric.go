package migration

import "database/sql"

func CreateCounterTable(db *sql.DB) error {
	query := `
create table counter_metric
(
    id        integer not null
        constraint counter_metric_pk
            primary key,
    type      varchar(20),
    value     integer
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
    type      varchar(20),
    value     precision double
);
`

	_, err := db.Exec(query)

	return err
}
