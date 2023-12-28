// Copyright © 2019 Binance
//
// This file is part of Binance. The full Binance copyright notice, including
// terms governing use, modification, and redistribution, is contained in the
// file LICENSE at the root of the source code distribution tree.

package keygen

import (
	"errors"
	"fmt"
	"math/big"

	"hwallet/tss/common"
	cmt "hwallet/tss/crypto/commitments"
	"hwallet/tss/crypto/vss"
	"hwallet/tss/tss"
)

// Implements Party
// Implements Stringer
var _ tss.Party = (*LocalParty)(nil)
var _ fmt.Stringer = (*LocalParty)(nil)

type (
	LocalParty struct {
		*tss.BaseParty
		params *tss.Parameters

		temp localTempData
		data LocalPartySaveData
	}

	localMessageStore struct {
		kgRound1Messages,
		kgRound2Message1s,
		kgRound2Message2s,
		kgRound3Messages,
		outMessages []tss.ParsedMessage
	}

	localTempData struct {
		localMessageStore

		// temp data (thrown away after keygen)
		ui            *big.Int // used for tests
		KGCs          []cmt.HashCommitment
		vs            vss.Vs
		shares        vss.Shares
		deCommitPolyG cmt.HashDeCommitment
	}
)

// Exported, used in `tss` client
func NewLocalParty(
	params *tss.Parameters,
	optionalPreParams ...LocalPreParams,
) tss.Party {
	partyCount := params.PartyCount()
	data := NewLocalPartySaveData(partyCount)
	// when `optionalPreParams` is provided we'll use the pre-computed primes instead of generating them from scratch
	if 0 < len(optionalPreParams) {
		if 1 < len(optionalPreParams) {
			panic(errors.New("keygen.NewLocalParty expected 0 or 1 item in `optionalPreParams`"))
		}
		if !optionalPreParams[0].ValidateWithProof() {
			panic(errors.New("`optionalPreParams` failed to validate; it might have been generated with an older version of tss-lib"))
		}
		//data.LocalPreParams = optionalPreParams[0]
	}
	p := &LocalParty{
		BaseParty: new(tss.BaseParty),
		params:    params,
		temp:      localTempData{},
		data:      data,
	}
	// msgs init
	p.temp.kgRound1Messages = make([]tss.ParsedMessage, partyCount)
	p.temp.kgRound2Message1s = make([]tss.ParsedMessage, partyCount)
	p.temp.kgRound2Message2s = make([]tss.ParsedMessage, partyCount)
	p.temp.kgRound3Messages = make([]tss.ParsedMessage, partyCount)
	// temp data init
	p.temp.KGCs = make([]cmt.HashCommitment, partyCount)
	return p
}

func (p *LocalParty) FirstRound() tss.Round {
	return newRound1(p.params, &p.data, &p.temp)
}

func (p *LocalParty) Start() (int, *tss.Error) {
	return tss.BaseStart(p, TaskName)
}

func (p *LocalParty) Update(msg tss.ParsedMessage) (roundNumber int, advanceRound bool, err *tss.Error) {
	roundNumber, advanceRound, err = tss.BaseUpdate(p, msg, TaskName)
	return
}

func (p *LocalParty) UpdateFromBytes(wireBytes []byte, from *tss.PartyID, isBroadcast bool) (int, bool, *tss.Error) {
	msg, err := tss.ParseWireMessage(wireBytes, from, isBroadcast)
	if err != nil {
		return -1, false, p.WrapError(err)
	}
	return p.Update(msg)
}

func (p *LocalParty) ValidateMessage(msg tss.ParsedMessage) (bool, *tss.Error) {
	if ok, err := p.BaseParty.ValidateMessage(msg); !ok || err != nil {
		return ok, err
	}
	// check that the message's "from index" will fit into the array
	if maxFromIdx := p.params.PartyCount() - 1; maxFromIdx < msg.GetFrom().Index {
		return false, p.WrapError(fmt.Errorf("received msg with a sender index too great (%d <= %d)",
			p.params.PartyCount(), msg.GetFrom().Index), msg.GetFrom())
	}
	return true, nil
}

func (p *LocalParty) GetOutMessage(index int) (tss.ParsedMessage, *tss.Error) {
	msg := p.temp.outMessages[index]
	return msg, nil
}

func (p *LocalParty) StoreMessage(msg tss.ParsedMessage) (bool, int, *tss.Error) {
	round := -1
	// ValidateBasic is cheap; double-check the message here in case the public StoreMessage was called externally
	if ok, err := p.ValidateMessage(msg); !ok || err != nil {
		return ok, round, err
	}
	fromPIdx := msg.GetFrom().Index

	// switch/case is necessary to store any messages beyond current round
	// this does not handle message replays. we expect the caller to apply replay and spoofing protection.
	switch msg.Content().(type) {
	case *KGRound1Message:
		p.temp.kgRound1Messages[fromPIdx] = msg
		round = 1
	case *KGRound2Message1:
		p.temp.kgRound2Message1s[fromPIdx] = msg
		round = 21
	case *KGRound2Message2:
		p.temp.kgRound2Message2s[fromPIdx] = msg
		round = 22
	case *KGRound3Message:
		p.temp.kgRound3Messages[fromPIdx] = msg
		round = 3
	default: // unrecognised message, just ignore!
		common.Logger.Warningf("unrecognised message ignored: %v", msg)
		return false, round, nil
	}
	return true, round, nil
}

// recovers a party's original index in the set of parties during keygen
func (save LocalPartySaveData) OriginalIndex() (int, error) {
	index := -1
	ki := save.ShareID
	for j, kj := range save.Ks {
		if kj.Cmp(ki) != 0 {
			continue
		}
		index = j
		break
	}
	if index < 0 {
		return -1, errors.New("a party index could not be recovered from Ks")
	}
	return index, nil
}

func (p *LocalParty) PartyID() *tss.PartyID {
	return p.params.PartyID()
}

func (p *LocalParty) String() string {
	return fmt.Sprintf("id: %s, %s", p.PartyID(), p.BaseParty.String())
}

func (p *LocalParty) SetNewSecretX(secretx *big.Int) {
	p.temp.ui = secretx
}

func (p *LocalParty) GetSavedata() LocalPartySaveData {
	return p.data
}
