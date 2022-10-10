-- DELETE visualization_variable_mapping entries when a visualization is deleted.

-- Drop the foreign key reference so we can add it back with an
-- ON DELETE CASCADE clause.

ALTER TABLE visualization_variable_mapping 
DROP CONSTRAINT visualization_variable_mapping_visualization_id_fkey,
ADD CONSTRAINT visualization_variable_mapping_visualization_id_fkey
    FOREIGN KEY (visualization_id)
    REFERENCES visualization(id)
    ON DELETE CASCADE;