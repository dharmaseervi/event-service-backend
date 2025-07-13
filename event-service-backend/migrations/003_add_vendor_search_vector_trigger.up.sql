-- Create search_vector column
ALTER TABLE vendors ADD COLUMN IF NOT EXISTS search_vector tsvector;

-- Create the function
CREATE OR REPLACE FUNCTION update_vendor_search_vector() RETURNS trigger AS $$
BEGIN
  NEW.search_vector :=
    to_tsvector('english', coalesce(NEW.title, '') || ' ' || coalesce(NEW.description, '') || ' ' || coalesce(NEW.category, '') || ' ' || coalesce(NEW.location, ''));
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger
CREATE TRIGGER vendor_search_vector_trigger
BEFORE INSERT OR UPDATE ON vendors
FOR EACH ROW
EXECUTE FUNCTION update_vendor_search_vector();
