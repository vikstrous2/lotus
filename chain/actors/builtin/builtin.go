package builtin

import (
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/cbor"

	"github.com/filecoin-project/lotus/chain/actors/adt"
	"github.com/filecoin-project/lotus/chain/types"

	builtin0 "github.com/filecoin-project/specs-actors/actors/builtin"
	miner0 "github.com/filecoin-project/specs-actors/actors/builtin/miner"
	proof0 "github.com/filecoin-project/specs-actors/actors/runtime/proof"
	smoothing0 "github.com/filecoin-project/specs-actors/actors/util/smoothing"
	builtin2 "github.com/filecoin-project/specs-actors/v2/actors/builtin"
	smoothing2 "github.com/filecoin-project/specs-actors/v2/actors/util/smoothing"
)

// TODO: Why does actors have 2 different versions of this?
type SectorInfo = proof0.SectorInfo
type PoStProof = proof0.PoStProof
type FilterEstimate = smoothing0.FilterEstimate

func FromV0FilterEstimate(v0 smoothing0.FilterEstimate) FilterEstimate {
	return (FilterEstimate)(v0)
}

// Doesn't change between actors v0 and v1
func QAPowerForWeight(size abi.SectorSize, duration abi.ChainEpoch, dealWeight, verifiedWeight abi.DealWeight) abi.StoragePower {
	return miner0.QAPowerForWeight(size, duration, dealWeight, verifiedWeight)
}

func FromV2FilterEstimate(v1 smoothing2.FilterEstimate) FilterEstimate {
	return (FilterEstimate)(v1)
}

type ActorStateLoader func(store adt.Store, root cid.Cid) (cbor.Marshaler, error)

var ActorStateLoaders = make(map[cid.Cid]ActorStateLoader)

func RegisterActorState(code cid.Cid, loader ActorStateLoader) {
	ActorStateLoaders[code] = loader
}

func Load(store adt.Store, act *types.Actor) (cbor.Marshaler, error) {
	loader, found := ActorStateLoaders[act.Code]
	if !found {
		return nil, xerrors.Errorf("unknown actor code %s", act.Code)
	}
	return loader(store, act.Head)
}

func ActorNameByCode(c cid.Cid) string {
	switch {
	case builtin0.IsBuiltinActor(c):
		return builtin0.ActorNameByCode(c)
	case builtin2.IsBuiltinActor(c):
		return builtin2.ActorNameByCode(c)
	default:
		return "<unknown>"
	}
}

func IsBuiltinActor(c cid.Cid) bool {
	return builtin0.IsBuiltinActor(c) || builtin2.IsBuiltinActor(c)
}