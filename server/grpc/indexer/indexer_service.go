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
	if chain.Path == "neutron" {
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
		if chain.Path == "neutron" {
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

func (i IndexerService) GetGovernanceProposalStati(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[indexerpb.GetGovernanceProposalStatiResponse], error) {
	chains := i.chainManager.QueryEnabled(ctx)
	var pbChains []*indexerpb.ChainInfo
	for _, chain := range chains {
		if chain.Path == "neutron" {
			continue
		}
		pbChains = append(pbChains, &indexerpb.ChainInfo{
			Id:        uint64(chain.ID),
			Name:      chain.Name,
			Path:      chain.Path,
			RpcUrl:    getRpcUrl(chain),
			Proposals: []*indexerpb.GovernanceProposal{},
		})
		proposals := i.chainManager.QueryProposals(ctx, chain)
		for _, proposal := range proposals {
			pbChains[len(pbChains)-1].Proposals = append(pbChains[len(pbChains)-1].Proposals, &indexerpb.GovernanceProposal{
				ProposalId: proposal.ProposalID,
				Status:     queryevent.ProposalStatus(queryevent.ProposalStatus_value[proposal.Status.String()]),
			})
		}
	}
	return connect.NewResponse(&indexerpb.GetGovernanceProposalStatiResponse{Chains: pbChains}), nil
}
