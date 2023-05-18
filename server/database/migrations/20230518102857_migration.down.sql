-- reverse: create index "chains_path_key" to table: "chains"
DROP INDEX "chains_path_key";
-- reverse: drop index "chains_chain_id_key" from table: "chains"
CREATE UNIQUE INDEX "chains_chain_id_key" ON "chains" ("chain_id");
