package shared

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrStateMetadataFileMissing is returned when the metadata file is missing.
var ErrStateMetadataFileMissing = errors.New("metadata file is missing")

// PostMetadata is the data associated with the PoST init procedure, persisted in the datadir next to the init files.
type PostMetadata struct {
	Version int `json:",omitempty"`

	NodeId          []byte
	CommitmentAtxId []byte

	LabelsPerUnit uint64
	NumUnits      uint32
	MaxFileSize   uint64
	Nonce         *uint64    `json:",omitempty"`
	NonceValue    NonceValue `json:",omitempty"`
	LastPosition  *uint64    `json:",omitempty"`
}

type NonceValue []byte

func (n NonceValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(n))
}

func (n *NonceValue) UnmarshalJSON(data []byte) (err error) {
	var hexString string
	if err = json.Unmarshal(data, &hexString); err != nil {
		return
	}
	*n, err = hex.DecodeString(hexString)
	return
}

const MetadataFileName = "postdata_metadata.json"

func SaveMetadata(dir string, v *PostMetadata) error {
	err := os.MkdirAll(dir, OwnerReadWriteExec)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("dir creation failure: %w", err)
	}

	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("serialization failure: %w", err)
	}

	err = os.WriteFile(filepath.Join(dir, MetadataFileName), data, OwnerReadWrite)
	if err != nil {
		return fmt.Errorf("write to disk failure: %w", err)
	}

	return nil
}

func LoadMetadata(dir string) (*PostMetadata, error) {
	filename := filepath.Join(dir, MetadataFileName)
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrStateMetadataFileMissing
		}
		return nil, fmt.Errorf("read file failure: %w", err)
	}

	metadata := PostMetadata{}
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}
