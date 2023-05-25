// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/loomi-labs/star-scope/ent/chain"
	"github.com/loomi-labs/star-scope/ent/contractproposal"
	"github.com/loomi-labs/star-scope/ent/event"
	"github.com/loomi-labs/star-scope/ent/eventlistener"
	"github.com/loomi-labs/star-scope/ent/proposal"
	"github.com/loomi-labs/star-scope/ent/schema"
	"github.com/loomi-labs/star-scope/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	chainMixin := schema.Chain{}.Mixin()
	chainMixinFields0 := chainMixin[0].Fields()
	_ = chainMixinFields0
	chainFields := schema.Chain{}.Fields()
	_ = chainFields
	// chainDescCreateTime is the schema descriptor for create_time field.
	chainDescCreateTime := chainMixinFields0[0].Descriptor()
	// chain.DefaultCreateTime holds the default value on creation for the create_time field.
	chain.DefaultCreateTime = chainDescCreateTime.Default.(func() time.Time)
	// chainDescUpdateTime is the schema descriptor for update_time field.
	chainDescUpdateTime := chainMixinFields0[1].Descriptor()
	// chain.DefaultUpdateTime holds the default value on creation for the update_time field.
	chain.DefaultUpdateTime = chainDescUpdateTime.Default.(func() time.Time)
	// chain.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	chain.UpdateDefaultUpdateTime = chainDescUpdateTime.UpdateDefault.(func() time.Time)
	// chainDescIndexingHeight is the schema descriptor for indexing_height field.
	chainDescIndexingHeight := chainFields[6].Descriptor()
	// chain.DefaultIndexingHeight holds the default value on creation for the indexing_height field.
	chain.DefaultIndexingHeight = chainDescIndexingHeight.Default.(uint64)
	// chainDescHasCustomIndexer is the schema descriptor for has_custom_indexer field.
	chainDescHasCustomIndexer := chainFields[7].Descriptor()
	// chain.DefaultHasCustomIndexer holds the default value on creation for the has_custom_indexer field.
	chain.DefaultHasCustomIndexer = chainDescHasCustomIndexer.Default.(bool)
	// chainDescHandledMessageTypes is the schema descriptor for handled_message_types field.
	chainDescHandledMessageTypes := chainFields[8].Descriptor()
	// chain.DefaultHandledMessageTypes holds the default value on creation for the handled_message_types field.
	chain.DefaultHandledMessageTypes = chainDescHandledMessageTypes.Default.(string)
	// chainDescUnhandledMessageTypes is the schema descriptor for unhandled_message_types field.
	chainDescUnhandledMessageTypes := chainFields[9].Descriptor()
	// chain.DefaultUnhandledMessageTypes holds the default value on creation for the unhandled_message_types field.
	chain.DefaultUnhandledMessageTypes = chainDescUnhandledMessageTypes.Default.(string)
	// chainDescIsEnabled is the schema descriptor for is_enabled field.
	chainDescIsEnabled := chainFields[10].Descriptor()
	// chain.DefaultIsEnabled holds the default value on creation for the is_enabled field.
	chain.DefaultIsEnabled = chainDescIsEnabled.Default.(bool)
	contractproposalMixin := schema.ContractProposal{}.Mixin()
	contractproposalMixinFields0 := contractproposalMixin[0].Fields()
	_ = contractproposalMixinFields0
	contractproposalFields := schema.ContractProposal{}.Fields()
	_ = contractproposalFields
	// contractproposalDescCreateTime is the schema descriptor for create_time field.
	contractproposalDescCreateTime := contractproposalMixinFields0[0].Descriptor()
	// contractproposal.DefaultCreateTime holds the default value on creation for the create_time field.
	contractproposal.DefaultCreateTime = contractproposalDescCreateTime.Default.(func() time.Time)
	// contractproposalDescUpdateTime is the schema descriptor for update_time field.
	contractproposalDescUpdateTime := contractproposalMixinFields0[1].Descriptor()
	// contractproposal.DefaultUpdateTime holds the default value on creation for the update_time field.
	contractproposal.DefaultUpdateTime = contractproposalDescUpdateTime.Default.(func() time.Time)
	// contractproposal.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	contractproposal.UpdateDefaultUpdateTime = contractproposalDescUpdateTime.UpdateDefault.(func() time.Time)
	eventMixin := schema.Event{}.Mixin()
	eventMixinFields0 := eventMixin[0].Fields()
	_ = eventMixinFields0
	eventFields := schema.Event{}.Fields()
	_ = eventFields
	// eventDescCreateTime is the schema descriptor for create_time field.
	eventDescCreateTime := eventMixinFields0[0].Descriptor()
	// event.DefaultCreateTime holds the default value on creation for the create_time field.
	event.DefaultCreateTime = eventDescCreateTime.Default.(func() time.Time)
	// eventDescUpdateTime is the schema descriptor for update_time field.
	eventDescUpdateTime := eventMixinFields0[1].Descriptor()
	// event.DefaultUpdateTime holds the default value on creation for the update_time field.
	event.DefaultUpdateTime = eventDescUpdateTime.Default.(func() time.Time)
	// event.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	event.UpdateDefaultUpdateTime = eventDescUpdateTime.UpdateDefault.(func() time.Time)
	// eventDescNotifyTime is the schema descriptor for notify_time field.
	eventDescNotifyTime := eventFields[2].Descriptor()
	// event.DefaultNotifyTime holds the default value on creation for the notify_time field.
	event.DefaultNotifyTime = eventDescNotifyTime.Default.(time.Time)
	eventlistenerMixin := schema.EventListener{}.Mixin()
	eventlistenerMixinFields0 := eventlistenerMixin[0].Fields()
	_ = eventlistenerMixinFields0
	eventlistenerFields := schema.EventListener{}.Fields()
	_ = eventlistenerFields
	// eventlistenerDescCreateTime is the schema descriptor for create_time field.
	eventlistenerDescCreateTime := eventlistenerMixinFields0[0].Descriptor()
	// eventlistener.DefaultCreateTime holds the default value on creation for the create_time field.
	eventlistener.DefaultCreateTime = eventlistenerDescCreateTime.Default.(func() time.Time)
	// eventlistenerDescUpdateTime is the schema descriptor for update_time field.
	eventlistenerDescUpdateTime := eventlistenerMixinFields0[1].Descriptor()
	// eventlistener.DefaultUpdateTime holds the default value on creation for the update_time field.
	eventlistener.DefaultUpdateTime = eventlistenerDescUpdateTime.Default.(func() time.Time)
	// eventlistener.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	eventlistener.UpdateDefaultUpdateTime = eventlistenerDescUpdateTime.UpdateDefault.(func() time.Time)
	proposalMixin := schema.Proposal{}.Mixin()
	proposalMixinFields0 := proposalMixin[0].Fields()
	_ = proposalMixinFields0
	proposalFields := schema.Proposal{}.Fields()
	_ = proposalFields
	// proposalDescCreateTime is the schema descriptor for create_time field.
	proposalDescCreateTime := proposalMixinFields0[0].Descriptor()
	// proposal.DefaultCreateTime holds the default value on creation for the create_time field.
	proposal.DefaultCreateTime = proposalDescCreateTime.Default.(func() time.Time)
	// proposalDescUpdateTime is the schema descriptor for update_time field.
	proposalDescUpdateTime := proposalMixinFields0[1].Descriptor()
	// proposal.DefaultUpdateTime holds the default value on creation for the update_time field.
	proposal.DefaultUpdateTime = proposalDescUpdateTime.Default.(func() time.Time)
	// proposal.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	proposal.UpdateDefaultUpdateTime = proposalDescUpdateTime.UpdateDefault.(func() time.Time)
	userMixin := schema.User{}.Mixin()
	userMixinFields0 := userMixin[0].Fields()
	_ = userMixinFields0
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreateTime is the schema descriptor for create_time field.
	userDescCreateTime := userMixinFields0[0].Descriptor()
	// user.DefaultCreateTime holds the default value on creation for the create_time field.
	user.DefaultCreateTime = userDescCreateTime.Default.(func() time.Time)
	// userDescUpdateTime is the schema descriptor for update_time field.
	userDescUpdateTime := userMixinFields0[1].Descriptor()
	// user.DefaultUpdateTime holds the default value on creation for the update_time field.
	user.DefaultUpdateTime = userDescUpdateTime.Default.(func() time.Time)
	// user.UpdateDefaultUpdateTime holds the default value on update for the update_time field.
	user.UpdateDefaultUpdateTime = userDescUpdateTime.UpdateDefault.(func() time.Time)
	// userDescWalletAddress is the schema descriptor for wallet_address field.
	userDescWalletAddress := userFields[1].Descriptor()
	// user.WalletAddressValidator is a validator for the "wallet_address" field. It is called by the builders before save.
	user.WalletAddressValidator = userDescWalletAddress.Validators[0].(func(string) error)
}
