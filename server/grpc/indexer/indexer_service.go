package indexer

import (
	"context"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/loomi-labs/star-scope/database"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/grpc/indexer/indexerpb"
	"github.com/loomi-labs/star-scope/grpc/indexer/indexerpb/indexerpbconnect"
	"github.com/loomi-labs/star-scope/queryevent"
	"github.com/shifty11/go-logger/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

//goland:noinspection GoNameStartsWithPackageName
type IndexerService struct {
	indexerpbconnect.UnimplementedIndexerServiceHandler
	chainManager *database.ChainManager
}

func NewIndexerServiceHandler(dbManagers *database.DbManagers) indexerpbconnect.IndexerServiceHandler {
	return &IndexerService{
		chainManager: dbManagers.ChainManager,
	}
}

func getRpcUrl(chain *ent.Chain) string {
	if chain.Path == "neutron-pion" {
		return "https://rest-palvus.pion-1.ntrn.tech"
	}
	return fmt.Sprintf("https://rest.cosmos.directory/%s", chain.Path)
}

func (i IndexerService) GetIndexingChains(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[indexerpb.GetIndexingChainsResponse], error) {
	chains := i.chainManager.QueryEnabled(ctx)
	pbChains := make([]*indexerpb.IndexingChain, len(chains))
	for ix, chain := range chains {
		pbChains[ix] = &indexerpb.IndexingChain{
			Id:                    uint64(chain.ID),
			Name:                  chain.Name,
			Path:                  chain.Path,
			RpcUrl:                getRpcUrl(chain),
			IndexingHeight:        chain.IndexingHeight,
			UnhandledMessageTypes: strings.Split(chain.UnhandledMessageTypes, ","),
			HasCustomIndexer:      chain.HasCustomIndexer,
		}
		if chain.Path == "neutron-pion" {
			pbChains[ix].RpcUrl = "https://rest-palvus.pion-1.ntrn.tech"
		}
	}

	return connect.NewResponse(&indexerpb.GetIndexingChainsResponse{Chains: pbChains}), nil
}

func (i IndexerService) UpdateIndexingChains(ctx context.Context, request *connect.Request[indexerpb.UpdateIndexingChainsRequest]) (*connect.Response[emptypb.Empty], error) {
	for _, chain := range request.Msg.GetChains() {
		err := i.chainManager.UpdateIndexStatus(ctx, int(chain.GetId()), chain.GetIndexingHeight(), chain.GetHandledMessageTypes(), chain.GetUnhandledMessageTypes())
		if err != nil {
			log.Sugar.Errorf("error while updating chain: %v", err)
		}
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (i IndexerService) shouldSendChain(chainPath string, chainPaths []string) bool {
	if len(chainPaths) == 0 {
		return true
	}
	for _, path := range chainPaths {
		if path == chainPath {
			return true
		}
	}
	return false
}

func (i IndexerService) GetGovernanceProposalStati(ctx context.Context, request *connect.Request[indexerpb.GetGovernanceProposalStatiRequest]) (*connect.Response[indexerpb.GetGovernanceProposalStatiResponse], error) {
	chains := i.chainManager.QueryEnabled(ctx)
	var pbChains []*indexerpb.ChainInfo
	for _, chain := range chains {
		if !i.shouldSendChain(chain.Path, request.Msg.GetChainPaths()) {
			continue
		}
		pbChains = append(pbChains, &indexerpb.ChainInfo{
			Id:                uint64(chain.ID),
			Name:              chain.Name,
			Path:              chain.Path,
			RpcUrl:            getRpcUrl(chain),
			Proposals:         []*indexerpb.GovernanceProposal{},
			ContractProposals: []*indexerpb.ContractGovernanceProposal{},
		})
		proposals := i.chainManager.QueryProposals(ctx, chain)
		for _, proposal := range proposals {
			pbChains[len(pbChains)-1].Proposals = append(pbChains[len(pbChains)-1].Proposals, &indexerpb.GovernanceProposal{
				ProposalId: proposal.ProposalID,
				Status:     queryevent.ProposalStatus(queryevent.ProposalStatus_value[proposal.Status.String()]),
			})
		}
		contractProposals := i.chainManager.QueryContractProposals(ctx, chain)
		for _, proposal := range contractProposals {
			pbChains[len(pbChains)-1].ContractProposals = append(pbChains[len(pbChains)-1].ContractProposals, &indexerpb.ContractGovernanceProposal{
				ProposalId:      proposal.ProposalID,
				Status:          queryevent.ContractProposalStatus(queryevent.ContractProposalStatus_value[proposal.Status.String()]),
				ContractAddress: proposal.ContractAddress,
			})
		}
	}
	return connect.NewResponse(&indexerpb.GetGovernanceProposalStatiResponse{Chains: pbChains}), nil
}
