-- reverse: create index "validator_operator_address_chain_validators" to table: "validators"
DROP INDEX "validator_operator_address_chain_validators";
-- reverse: create index "validator_operator_address" to table: "validators"
DROP INDEX "validator_operator_address";
-- reverse: create index "validator_moniker_operator_address_chain_validators" to table: "validators"
DROP INDEX "validator_moniker_operator_address_chain_validators";
-- reverse: create index "validator_moniker_address_chain_validators" to table: "validators"
DROP INDEX "validator_moniker_address_chain_validators";
-- reverse: create index "validator_moniker" to table: "validators"
DROP INDEX "validator_moniker";
-- reverse: create index "validator_address_chain_validators" to table: "validators"
DROP INDEX "validator_address_chain_validators";
-- reverse: create index "validator_address" to table: "validators"
DROP INDEX "validator_address";
-- reverse: create "validators" table
DROP TABLE "validators";
