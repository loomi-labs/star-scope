package user

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/commchannel"
	"github.com/loomi-labs/star-scope/ent/usersetup"
	"github.com/loomi-labs/star-scope/grpc/types"
	"github.com/loomi-labs/star-scope/grpc/user/userpb"
	"github.com/loomi-labs/star-scope/grpc/user/userpb/userpbconnect"
	sf "github.com/sa-/slicefunk"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (u *UserSetupService) createStepResponse(ctx context.Context, setup *ent.UserSetup, requestedStep usersetup.Step, isComplete bool) *userpb.StepResponse {
	response := &userpb.StepResponse{}
	switch requestedStep {
	case usersetup.StepOne:
		response.Step = &userpb.StepResponse_One{One: &userpb.StepOneResponse{
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
		response.Step = &userpb.StepResponse_Two{Two: &userpb.StepTwoResponse{
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
		response.Step = &userpb.StepResponse_Three{Three: &userpb.StepThreeResponse{
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
		response.Step = &userpb.StepResponse_Four{Four: &userpb.StepFourResponse{
			NotifyFunding:           setup.NotifyFunding,
			NotifyStaking:           setup.NotifyStaking,
			NotifyGovNewProposal:    setup.NotifyGovNewProposal,
			NotifyGovVotingEnd:      setup.NotifyGovVotingEnd,
			NotifyGovVotingReminder: setup.NotifyGovVotingReminder,
			NotifyGovChainIds:       selectedIds,
			AvailableChains:         availableChains,
		}}
	case usersetup.StepFive:
		response.Step = &userpb.StepResponse_Five{Five: &userpb.StepFiveResponse{}}
	}
	response.NumSteps = u.getNumSteps(setup)
	response.IsComplete = isComplete
	return response
}

func (u *UserSetupService) getNumSteps(setup *ent.UserSetup) uint32 {
	if setup.IsValidator {
		return 5
	} else {
		return 4
	}
}

func (u *UserSetupService) GetStep(ctx context.Context, request *connect.Request[userpb.GetStepRequest]) (*connect.Response[userpb.StepResponse], error) {
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

	step := setup.Step
	if request.Msg != nil {
		switch request.Msg.GetStep() {
		case userpb.GetStepRequest_CURRENT:
			break
		case userpb.GetStepRequest_ONE:
			step = usersetup.StepOne
		case userpb.GetStepRequest_TWO:
			step = usersetup.StepTwo
		case userpb.GetStepRequest_THREE:
			step = usersetup.StepThree
		case userpb.GetStepRequest_FOUR:
			step = usersetup.StepFour
		case userpb.GetStepRequest_FIVE:
			step = usersetup.StepFive
		}
	}
	response := u.createStepResponse(ctx, setup, step, user.IsSetupComplete)

	return connect.NewResponse(response), nil
}

func isFinishStepRequestValid(request *connect.Request[userpb.FinishStepRequest]) bool {
	if request.Msg == nil || request.Msg.Step == nil {
		return false
	}
	switch request.Msg.Step.(type) {
	case *userpb.FinishStepRequest_One:
		return request.Msg.GetOne() != nil
	case *userpb.FinishStepRequest_Two:
		return request.Msg.GetTwo() != nil
	case *userpb.FinishStepRequest_Three:
		return request.Msg.GetThree() != nil
	case *userpb.FinishStepRequest_Four:
		return request.Msg.GetFour() != nil
	case *userpb.FinishStepRequest_Five:
		return request.Msg.GetFive() != nil
	}
	return false
}

func (u *UserSetupService) FinishStep(ctx context.Context, request *connect.Request[userpb.FinishStepRequest]) (*connect.Response[userpb.StepResponse], error) {
	user, ok := ctx.Value(common.ContextKeyUser).(*ent.User)
	if !ok {
		log.Sugar.Error("invalid user")
		return nil, types.UserNotFoundErr
	}

	if user.IsSetupComplete {
		log.Sugar.Errorf("user setup already completed", "user", user)
		return nil, status.Error(codes.InvalidArgument, "Setup already completed")
	}

	if !isFinishStepRequestValid(request) {
		log.Sugar.Errorf("invalid request for finish nextStep: %v", request)
		return nil, types.InvalidArgumentErr
	}

	setup, err := u.userManager.QuerySetup(ctx, user)
	if err != nil {
		log.Sugar.Errorf("failed to query setup", "error", err)
		return nil, types.UnknownErr
	}

	var isCompleted = false
	var updateQuery *ent.UserSetupUpdateOne
	switch request.Msg.Step.(type) {
	case *userpb.FinishStepRequest_One:
		var step = usersetup.StepThree
		if request.Msg.GetOne().GetIsValidator() {
			step = usersetup.StepTwo
		}
		updateQuery = setup.
			Update().
			SetIsValidator(request.Msg.GetOne().GetIsValidator()).
			SetStep(step)
	case *userpb.FinishStepRequest_Two:
		step := usersetup.StepThree
		if !request.Msg.GetGoToNextStep() {
			step = usersetup.StepOne
		}
		validatorIds := sf.Map(request.Msg.GetTwo().GetValidatorIds(), func(id int64) int { return int(id) })
		updateQuery = setup.
			Update().
			ClearSelectedValidators().
			AddSelectedValidatorIDs(validatorIds...).
			SetStep(step)
	case *userpb.FinishStepRequest_Three:
		step := usersetup.StepFour
		if !request.Msg.GetGoToNextStep() {
			if setup.IsValidator {
				step = usersetup.StepTwo
			} else {
				step = usersetup.StepOne
			}
		}
		var chainIds []int
		var chains = u.chainManager.QueryEnabled(ctx)
		for _, address := range request.Msg.GetThree().GetWalletAddresses() {
			for _, chain := range chains {
				if common.IsBech32AddressFromChain(address, chain.Bech32Prefix) {
					chainIds = append(chainIds, chain.ID)
					break
				}
			}
		}
		updateQuery = setup.
			Update().
			SetWalletAddresses(request.Msg.GetThree().GetWalletAddresses()).
			ClearSelectedChains().
			AddSelectedChainIDs(sf.Unique(chainIds)...).
			SetStep(step)
	case *userpb.FinishStepRequest_Four:
		step := usersetup.StepFive
		if !request.Msg.GetGoToNextStep() {
			step = usersetup.StepThree
		}
		notifyGovChainIds := sf.Map(request.Msg.GetFour().GetNotifyGovChainIds(), func(id int64) int { return int(id) })
		updateQuery = setup.
			Update().
			SetNotifyFunding(request.Msg.GetFour().GetNotifyFunding()).
			SetNotifyStaking(request.Msg.GetFour().GetNotifyStaking()).
			SetNotifyGovNewProposal(request.Msg.GetFour().GetNotifyGovNewProposal()).
			SetNotifyGovVotingEnd(request.Msg.GetFour().GetNotifyGovVotingEnd()).
			SetNotifyGovVotingReminder(request.Msg.GetFour().GetNotifyGovVotingReminder()).
			ClearSelectedChains().
			AddSelectedChainIDs(notifyGovChainIds...).
			SetStep(step)
	case *userpb.FinishStepRequest_Five:
		step := usersetup.StepFive
		if !request.Msg.GetGoToNextStep() {
			step = usersetup.StepFour
		} else {
			isCompleted = true
		}
		switch request.Msg.GetFive().GetChannel().(type) {
		case *userpb.StepFiveRequest_Webapp:
		case *userpb.StepFiveRequest_Telegram:
			if user.TelegramUserID == 0 {
				log.Sugar.Errorf("invalid telegram user id: %v", user.TelegramUserID)
				return nil, types.InvalidArgumentErr
			}
			t := commchannel.TypeTelegram
			channels, err := u.userManager.QueryCommChannels(ctx, user, &t)
			if err != nil {
				log.Sugar.Errorf("failed to query comm channels: %v", err)
				return nil, types.UnknownErr
			}
			var found = false
			for _, channel := range channels {
				if channel.TelegramChatID == request.Msg.GetFive().GetTelegram().GetChatId() {
					found = true
					break
				}
			}
			if !found {
				log.Sugar.Errorf("invalid telegram chat id: %v", request.Msg.GetFive().GetTelegram().GetChatId())
				return nil, types.InvalidArgumentErr
			}
		case *userpb.StepFiveRequest_Discord:
			if user.DiscordUserID == 0 {
				return nil, types.InvalidArgumentErr
			}
			t := commchannel.TypeDiscord
			channels, err := u.userManager.QueryCommChannels(ctx, user, &t)
			if err != nil {
				log.Sugar.Errorf("failed to query comm channels: %v", err)
				return nil, types.UnknownErr
			}
			var found = false
			for _, channel := range channels {
				if channel.DiscordChannelID == request.Msg.GetFive().GetDiscord().GetChannelId() {
					found = true
					break
				}
			}
			if !found {
				log.Sugar.Errorf("invalid discord channel id: %v", request.Msg.GetFive().GetDiscord().GetChannelId())
				return nil, types.InvalidArgumentErr
			}
		}
		updateQuery = setup.
			Update().
			SetStep(step)
	}
	var chains []*ent.Chain
	if isCompleted {
		chains = u.chainManager.QueryEnabled(ctx)
	}
	setup, err = u.userManager.UpdateSetup(ctx, user, updateQuery, isCompleted, chains)
	if err != nil {
		log.Sugar.Errorf("failed to update setup", "error", err)
		return nil, types.UnknownErr
	}

	response := u.createStepResponse(ctx, setup, setup.Step, isCompleted)

	return connect.NewResponse(response), nil
}

func (u *UserSetupService) ValidateWallet(ctx context.Context, request *connect.Request[userpb.ValidateWalletRequest]) (*connect.Response[userpb.ValidateWalletResponse], error) {
	err := common.ValidateBech32Address(request.Msg.GetAddress())
	isValid := err == nil

	var wallet = &userpb.Wallet{
		Address: request.Msg.GetAddress(),
		LogoUrl: "",
	}
	var isSupported = false
	for _, chain := range u.chainManager.QueryEnabled(ctx) {
		if common.IsBech32AddressFromChain(request.Msg.GetAddress(), chain.Bech32Prefix) {
			wallet.LogoUrl = chain.Image
			isSupported = true
			break
		}
	}

	return connect.NewResponse(&userpb.ValidateWalletResponse{
		IsValid:     isValid,
		IsSupported: isSupported,
		Wallet:      wallet,
	}), nil
}
