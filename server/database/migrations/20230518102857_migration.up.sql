-- drop index "chains_chain_id_key" from table: "chains"
DROP INDEX "chains_chain_id_key";
-- create index "chains_path_key" to table: "chains"
CREATE UNIQUE INDEX "chains_path_key" ON "chains" ("path");
