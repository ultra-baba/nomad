package structs

import (
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/hashicorp/nomad/plugins/csi"
)

// CSIVolumeMountOptions contains the mount options that should be provided when
// attaching and mounting a volume with the CSIVolumeAttachmentModeFilesystem
// attachment mode.
type CSIVolumeMountOptions struct {
	// Filesystem is the desired filesystem type that should be used by the volume
	// (e.g ext4, aufs, zfs). This field is optional.
	Filesystem string

	// MountFlags contain the mount options that should be used for the volume.
	// These may contain _sensitive_ data and should not be leaked to logs or
	// returned in debugging data.
	// The total size of this field must be under 4KiB.
	MountFlags []string
}

type ClientCSIControllerValidateVolumeRequest struct {
	PluginID string
	VolumeID string

	AttachmentMode structs.CSIVolumeAttachmentMode
	AccessMode     structs.CSIVolumeAccessMode

	NodeID string

	structs.QueryOptions
}

type ClientCSIControllerValidateVolumeResponse struct {
}

type ClientCSIControllerAttachVolumeRequest struct {
	PluginName string

	// The ID of the volume to be used on a node.
	// This field is REQUIRED.
	VolumeID string

	// The ID of the node in CSI Terms. This field is REQUIRED. This must match the NodeID that
	// is fingerprinted by the target node for this plugin name.
	CSINodeID string

	// AttachmentMode indicates how the volume should be attached and mounted into
	// a task.
	AttachmentMode structs.CSIVolumeAttachmentMode

	// AccessMode indicates the desired concurrent access model for the volume
	AccessMode structs.CSIVolumeAccessMode

	// MountOptions is an optional field that contains additional configuration
	// when providing an AttachmentMode of CSIVolumeAttachmentModeFilesystem
	MountOptions *CSIVolumeMountOptions

	// ReadOnly indicates that the volume will be used in a readonly fashion. This
	// only works when the Controller has the PublishReadonly capability.
	ReadOnly bool

	// NodeID is the target ID of the CSIController node that should recieve the
	// RPC
	NodeID string

	structs.QueryOptions
}

func (c *ClientCSIControllerAttachVolumeRequest) ToCSIRequest() *csi.ControllerPublishVolumeRequest {
	if c == nil {
		return &csi.ControllerPublishVolumeRequest{}
	}

	caps, err := csi.VolumeCapabilityFromStructs(c.AttachmentMode, c.AccessMode)
	if err != nil {
		return nil // TODO(dani): Return errors here
	}

	return &csi.ControllerPublishVolumeRequest{
		VolumeID:         c.VolumeID,
		NodeID:           c.CSINodeID,
		ReadOnly:         c.ReadOnly,
		VolumeCapability: caps,
	}
}

type ClientCSIControllerAttachVolumeResponse struct {
	// Opaque static publish properties of the volume. SP MAY use this
	// field to ensure subsequent `NodeStageVolume` or `NodePublishVolume`
	// calls calls have contextual information.
	// The contents of this field SHALL be opaque to nomad.
	// The contents of this field SHALL NOT be mutable.
	// The contents of this field SHALL be safe for the nomad to cache.
	// The contents of this field SHOULD NOT contain sensitive
	// information.
	// The contents of this field SHOULD NOT be used for uniquely
	// identifying a volume. The `volume_id` alone SHOULD be sufficient to
	// identify the volume.
	// This field is OPTIONAL and when present MUST be passed to
	// subsequent `NodeStageVolume` or `NodePublishVolume` calls
	PublishContext map[string]string
}

type ClientCSIControllerDetachVolumeRequest struct {
	PluginName string

	// The ID of the volume to be unpublished for the node
	// This field is REQUIRED.
	VolumeID string

	// The ID of the node. This field is REQUIRED. This must match the NodeID that
	// is fingerprinted by the target node for this plugin name.
	NodeID string
}

func (c *ClientCSIControllerDetachVolumeRequest) ToCSIRequest() *csi.ControllerUnpublishVolumeRequest {
	if c == nil {
		return &csi.ControllerUnpublishVolumeRequest{}
	}

	return &csi.ControllerUnpublishVolumeRequest{
		VolumeID: c.VolumeID,
		NodeID:   c.NodeID,
	}
}

type ClientCSIControllerDetachVolumeResponse struct{}