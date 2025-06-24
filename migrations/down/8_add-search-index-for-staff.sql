DROP INDEX IF EXISTS staff_search_idx;

ALTER TABLE staff.staff DROP COLUMN IF EXISTS search_vector;

DROP FUNCTION IF EXISTS staff.staff_search_trigger;
