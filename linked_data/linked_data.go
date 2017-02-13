package linked_data

import (
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/spec"
)

func ValidateModelId(modelId string) (Data, error) {
	tx, err := bigchain.GetTx(modelId)
	if err != nil {
		return nil, err
	}
	model := bigchain.GetTxData(tx)
	if err = ValidateModel(model); err != nil {
		return nil, err
	}
	return model, nil
}

func ValidateMusicId(musicId string) (Data, error) {
	tx, err := bigchain.GetTx(musicId)
	if err != nil {
		return nil, err
	}
	music := bigchain.GetTxData(tx)
	if err = ValidateMusic(music); err != nil {
		return nil, err
	}
	return music, nil
}

func ValidateRightId(rightId string) (Data, error) {
	tx, err := bigchain.GetTx(rightId)
	if err != nil {
		return nil, err
	}
	right := bigchain.GetTxData(tx)
	if err := ValidateRight(right); err != nil {
		return nil, err
	}
	return right, nil
}

func ValidateSignatureId(signatureId string) (Data, error) {
	tx, err := bigchain.GetTx(signatureId)
	if err != nil {
		return nil, err
	}
	signature := bigchain.GetTxData(tx)
	if err = ValidateSignature(signature); err != nil {
		return nil, err
	}
	return signature, nil
}

func ValidateModel(model Data) error {
	_type := spec.GetType(model)
	switch _type {
	case spec.ALBUM:
		return ValidateAlbum(model)
	case spec.TRACK:
		return ValidateTrack(model)
	case spec.SIGNATURE:
		return ValidateSignature(model)
	case spec.RIGHT:
		return ValidateRight(model)
	}
	return ErrorAppend(ErrInvalidType, _type)
}

func ValidateMusic(music Data) error {
	_type := spec.GetType(music)
	if _type == spec.ALBUM {
		return ValidateAlbum(music)
	}
	if _type == spec.TRACK {
		return ValidateTrack(music)
	}
	return ErrorAppend(ErrInvalidType, _type)
}

func ValidateAlbum(album Data) error {
	if !spec.ValidAlbum(album) {
		return ErrorAppend(ErrInvalidModel, spec.ALBUM)
	}
	artistId := spec.GetMusicArtist(album)
	tx, err := bigchain.GetTx(artistId)
	if err != nil {
		return err
	}
	artist := bigchain.GetTxData(tx)
	if !spec.ValidArtist(artist) {
		return ErrorAppend(ErrInvalidModel, spec.ARTIST)
	}
	publisherId := spec.GetMusicPublisher(album)
	tx, err = bigchain.GetTx(publisherId)
	if err != nil {
		return err
	}
	publisher := bigchain.GetTxData(tx)
	if !spec.ValidPublisher(publisher) {
		return ErrorAppend(ErrInvalidModel, spec.PUBLISHER)
	}
	return nil
}

func ValidateTrack(track Data) error {
	if !spec.ValidTrack(track) {
		return ErrorAppend(ErrInvalidModel, spec.TRACK)
	}
	artistId := spec.GetMusicArtist(track)
	tx, err := bigchain.GetTx(artistId)
	if err != nil {
		return err
	}
	artist := bigchain.GetTxData(tx)
	if !spec.ValidArtist(artist) {
		return ErrorAppend(ErrInvalidModel, spec.ARTIST)
	}
	publisherId := spec.GetMusicPublisher(track)
	if publisherId != "" {
		tx, err = bigchain.GetTx(publisherId)
		if err != nil {
			return err
		}
		publisher := bigchain.GetTxData(tx)
		if !spec.ValidPublisher(publisher) {
			return ErrorAppend(ErrInvalidModel, spec.PUBLISHER)
		}
		return nil
	}
	albumId := spec.GetTrackAlbum(track)
	tx, err = bigchain.GetTx(albumId)
	if err != nil {
		return err
	}
	album := bigchain.GetTxData(tx)
	return ValidateAlbum(album)
}

func ValidateSignature(signature Data) error {
	if !spec.ValidSignature(signature) {
		return ErrorAppend(ErrInvalidModel, spec.SIGNATURE)
	}
	modelId := spec.GetSignatureModel(signature)
	tx, err := bigchain.GetTx(modelId)
	if err != nil {
		return err
	}
	model := bigchain.GetTxData(tx)
	if err := ValidateModel(model); err != nil {
		return err
	}
	signerId := spec.GetSignatureSigner(signature)
	tx, err = bigchain.GetTx(signerId)
	if err != nil {
		return err
	}
	signer := bigchain.GetTxData(tx)
	if !spec.ValidAgent(signer) {
		return err
	}
	pub := spec.GetAgentPublicKey(signer)
	sig := spec.GetSignatureValue(signature)
	if !pub.Verify(MustMarshalJSON(model), sig) {
		return ErrInvalidSignature
	}
	return nil
}

func ValidateRight(right Data) error {
	if !spec.ValidRight(right) {
		return ErrorAppend(ErrInvalidModel, spec.RIGHT)
	}
	musicId := spec.GetRightMusic(right)
	tx, err := bigchain.GetTx(musicId)
	if err != nil {
		return err
	}
	music := bigchain.GetTxData(tx)
	if err = ValidateMusic(music); err != nil {
		return err
	}
	artistId := spec.GetMusicArtist(music)
	tx, err = bigchain.GetTx(artistId)
	if err != nil {
		return err
	}
	artist := bigchain.GetTxData(tx)
	if !spec.ValidArtist(artist) {
		return ErrorAppend(ErrInvalidModel, spec.ARTIST)
	}
	pub := spec.GetAgentPublicKey(artist)
	recipientId := spec.GetRightRecipient(right)
	tx, err = bigchain.GetTx(recipientId)
	if err != nil {
		return err
	}
	recipient := bigchain.GetTxData(tx)
	if !spec.ValidAgent(recipient) {
		return ErrorAppend(ErrInvalidModel, spec.GetType(recipient))
	}
	signature := spec.GetRightSignature(right)
	if err = ValidateSignature(signature); err != nil {
		return err
	}
	sig := spec.GetSignatureValue(signature)
	if !pub.Verify(MustMarshalJSON(music), sig) {
		return ErrInvalidSignature
	}
	return nil
}
