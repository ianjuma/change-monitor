package sa

import (
	"context"
	"strings"
	"time"

	"database/sql"

	log "github.com/sirupsen/logrus"
)

// createTriggers creates db triggers and functions defined at constant query value
func createTriggers(conninfo string) {
	const op = "createTriggers"

	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		log.WithFields(log.Fields{
			"op":  op,
			"err": err,
		}).Error("failed to open database connection")
		return
	}

	parts := strings.Split(query, "---")
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	for _, part := range parts {
		if _, err := db.ExecContext(ctx, strings.TrimSpace(part)); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				log.WithFields(log.Fields{
					"op":   op,
					"err":  err,
					"part": part,
				}).Error("error creating trigger?")
			}
		}
	}
	if err := db.Close(); err != nil {
		log.WithFields(log.Fields{
			"op":  op,
			"err": err,
		}).Error("failed to close database connection")
	}
}

// CREATE TRIGGER does not have a "replace" semantic, will produce errors when it already exists
// queries must be separated with "---"
const query = `
CREATE OR REPLACE FUNCTION fn_notify_product_active_change()
RETURNS trigger AS $psql$
BEGIN
  PERFORM pg_notify(
    'product_active_change_event',
    json_build_object(
      'operation', TG_OP,
      'old', row_to_json(OLD),
      'new', row_to_json(NEW)
    )::text
  );RETURN NEW;
END;$psql$ language plpgsql;
---

CREATE TRIGGER product_active_change_event
AFTER UPDATE OF active
ON "product"
FOR EACH ROW
EXECUTE PROCEDURE fn_notify_product_active_change();
---

CREATE OR REPLACE FUNCTION fn_notify_product_insert()
    RETURNS trigger AS $psql$
BEGIN
    PERFORM pg_notify(
            'product_insert_event',
            json_build_object(
                    'operation', TG_OP,
                    'new', row_to_json(NEW)
                )::text
        );RETURN NEW;
END;$psql$ language plpgsql;
---

CREATE TRIGGER product_insert_event
    AFTER INSERT
    ON "product"
    FOR EACH ROW
EXECUTE PROCEDURE fn_notify_product_insert();
`
