package user

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/usersetup"
	"github.com/loomi-labs/star-scope/grpc/types"
	"github.com/loomi-labs/star-scope/grpc/user/userpb"
	"github.com/loomi-labs/star-scope/grpc/user/userpb/userpbconnect"
	sf "github.com/sa-/slicefunk"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

//goland:noinspection GoNameStartsWithPackageName
type UserSetupService struct {
	userpbconnect.UnimplementedUserSetupServiceHandler
	userManager      *database.UserManager
	chainManager     *database.ChainManager
	validatorManager *database.ValidatorManager
}

func NewUserSetupServiceHandler(dbManagers *database.DbManagers) userpbconnect.UserSetupServiceHandler {
	return &UserSetupService{
		userManager:      dbManagers.UserManager,
		chainManager:     dbManagers.ChainManager,
		validatorManager: dbManagers.ValidatorManager,
	}
}

func (u *UserSetupService) createStepResponse(ctx context.Context, setup *ent.UserSetup) *userpb.StepResponse {
	response := &userpb.StepResponse{}
	switch setup.Step {
	case usersetup.StepOne:
		response.Step = &userpb.StepResponse_StepOne{StepOne: &userpb.StepOneResponse{
			IsValidator: setup.IsValidator,
		}}
	case usersetup.StepTwo:
		bundles := u.validatorManager.QueryActiveBundledByMoniker(ctx)
		availableValidators := sf.Map(bundles, func(bundle *database.ValidatorBundle) *userpb.Validator {
			ids := sf.Map(bundle.Validators, func(validator *ent.Validator) int64 { return int64(validator.ID) })
			return &userpb.Validator{
				Ids:     ids,
				Moniker: bundle.Moniker,
			}
		})
		selectedIds := sf.Map(setup.QuerySelectedValidators().IDsX(ctx), func(id int) int64 { return int64(id) })
		response.Step = &userpb.StepResponse_StepTwo{StepTwo: &userpb.StepTwoResponse{
			AvailableValidators:  availableValidators,
			SelectedValidatorIds: selectedIds,
		}}
	case usersetup.StepThree:
		var wallets []*userpb.Wallet
		var chains = u.chainManager.QueryAll(ctx)
		for _, wallet := range setup.WalletAddresses {
			var logoUrl string
			for _, chain := range chains {
				if common.IsBech32AddressFromChain(wallet, chain.Bech32Prefix) {
					logoUrl = chain.Image
					break
				}
			}
			wallets = append(wallets, &userpb.Wallet{
				Address: wallet,
				LogoUrl: logoUrl,
			})
		}
		response.Step = &userpb.StepResponse_StepThree{StepThree: &userpb.StepThreeResponse{
			Wallets: wallets,
		}}
	case usersetup.StepFour:
		enabled := u.chainManager.QueryEnabled(ctx)
		availableChains := sf.Map(enabled, func(chain *ent.Chain) *userpb.GovChain {
			return &userpb.GovChain{
				Id:      int64(chain.ID),
				Name:    chain.Name,
				LogoUrl: chain.Image,
			}
		})
		selectedIds := sf.Map(setup.QuerySelectedChains().IDsX(ctx), func(id int) int64 { return int64(id) })
		response.Step = &userpb.StepResponse_StepFour{StepFour: &userpb.StepFourResponse{
			NotifyFunding:           setup.NotifyFunding,
			NotifyStaking:           setup.NotifyStaking,
			NotifyGovNewProposal:    setup.NotifyGovNewProposal,
			NotifyGovVotingEnd:      setup.NotifyGovVotingEnd,
			NotifyGovVotingReminder: setup.NotifyGovVotingReminder,
			NotifyGovChainIds:       selectedIds,
			AvailableChains:         availableChains,
		}}
	case usersetup.StepFive:
		response.Step = &userpb.StepResponse_StepFive{StepFive: &userpb.StepFiveResponse{}}
	}

	return response
}

func (u *UserSetupService) GetCurrentStep(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[userpb.StepResponse], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	setup, err := u.userManager.QuerySetup(ctx, user)
	if err != nil {
		log.Sugar.Errorf("failed to query setup", "error", err)
		return nil, types.UnknownErr
	}

	response := u.createStepResponse(ctx, setup)

	return connect.NewResponse(response), nil
}

func isFinishStepRequestValid(request *connect.Request[userpb.FinishStepRequest]) bool {
	if request.Msg == nil || request.Msg.Step == nil {
		return false
	}
	switch request.Msg.Step.(type) {
	case *userpb.FinishStepRequest_StepOne:
		return request.Msg.GetStepOne() != nil
	case *userpb.FinishStepRequest_StepTwo:
		return request.Msg.GetStepTwo() != nil && request.Msg.GetStepTwo().ValidatorIds != nil
	case *userpb.FinishStepRequest_StepThree:
		return request.Msg.GetStepThree() != nil && request.Msg.GetStepThree().WalletAddresses != nil
	case *userpb.FinishStepRequest_StepFour:
		return request.Msg.GetStepFour() != nil && request.Msg.GetStepFour().NotifyGovChainIds != nil
	case *userpb.FinishStepRequest_StepFive:
		return request.Msg.GetStepFive() != nil
	}
	return false
}

func (u *UserSetupService) FinishStep(ctx context.Context, request *connect.Request[userpb.FinishStepRequest]) (*connect.Response[userpb.StepResponse], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	if !isFinishStepRequestValid(request) {
		log.Sugar.Errorf("invalid request for finish step: %v", request)
		return nil, types.InvalidArgumentErr
	}

	setup, err := u.userManager.QuerySetup(ctx, user)
	if err != nil {
		log.Sugar.Errorf("failed to query setup", "error", err)
		return nil, types.UnknownErr
	}

	var updateQuery *ent.UserSetupUpdateOne
	switch request.Msg.Step.(type) {
	case *userpb.FinishStepRequest_StepOne:
		var step = usersetup.StepThree
		if request.Msg.GetStepOne().GetIsValidator() {
			setup.Step = usersetup.StepTwo
		}
		updateQuery = setup.
			Update().
			SetIsValidator(request.Msg.GetStepOne().GetIsValidator()).
			SetStep(step)
	case *userpb.FinishStepRequest_StepTwo:
		validatorIds := sf.Map(request.Msg.GetStepTwo().GetValidatorIds(), func(id int64) int { return int(id) })
		updateQuery = setup.
			Update().
			AddSelectedValidatorIDs(validatorIds...)
	case *userpb.FinishStepRequest_StepThree:
		updateQuery = setup.
			Update().
			SetWalletAddresses(request.Msg.GetStepThree().GetWalletAddresses())
	case *userpb.FinishStepRequest_StepFour:
		notifyGovChainIds := sf.Map(request.Msg.GetStepFour().GetNotifyGovChainIds(), func(id int64) int { return int(id) })
		updateQuery = setup.
			Update().
			SetNotifyFunding(request.Msg.GetStepFour().GetNotifyFunding()).
			SetNotifyStaking(request.Msg.GetStepFour().GetNotifyStaking()).
			SetNotifyGovNewProposal(request.Msg.GetStepFour().GetNotifyGovNewProposal()).
			SetNotifyGovVotingEnd(request.Msg.GetStepFour().GetNotifyGovVotingEnd()).
			SetNotifyGovVotingReminder(request.Msg.GetStepFour().GetNotifyGovVotingReminder()).
			AddSelectedChainIDs(notifyGovChainIds...)
	case *userpb.FinishStepRequest_StepFive:
		// TODO: update user setup
	}
	setup, err = u.userManager.UpdateSetup(ctx, updateQuery)
	if err != nil {
		log.Sugar.Errorf("failed to update setup", "error", err)
		return nil, types.UnknownErr
	}

	response := u.createStepResponse(ctx, setup)

	return connect.NewResponse(response), nil
}
