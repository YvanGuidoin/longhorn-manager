package types

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/longhorn/longhorn-manager/util"
)

const (
	LonghornKindNode                = "Node"
	LonghornKindVolume              = "Volume"
	LonghornKindEngineImage         = "EngineImage"
	LonghornKindInstanceManager     = "InstanceManager"
	LonghornKindShareManager        = "ShareManager"
	LonghornKindBackingImage        = "BackingImage"
	LonghornKindBackingImageManager = "BackingImageManager"

	LonghornKindBackingImageDataSource = "BackingImageDataSource"

	CRDAPIVersionV1alpha1 = "longhorn.rancher.io/v1alpha1"
	CRDAPIVersionV1beta1  = "longhorn.io/v1beta1"
	CurrentCRDAPIVersion  = CRDAPIVersionV1beta1
)

const (
	DefaultAPIPort = 9500

	EngineBinaryDirectoryInContainer = "/engine-binaries/"
	EngineBinaryDirectoryOnHost      = "/var/lib/longhorn/engine-binaries/"
	ReplicaHostPrefix                = "/host"
	EngineBinaryName                 = "longhorn"

	BackingImagesManagerDirectory = "/backing-images/"
	BackingImageFileName          = "backing"

	LonghornNodeKey     = "longhornnode"
	LonghornDiskUUIDKey = "longhorndiskuuid"

	NodeCreateDefaultDiskLabelKey             = "node.longhorn.io/create-default-disk"
	NodeCreateDefaultDiskLabelValueTrue       = "true"
	NodeCreateDefaultDiskLabelValueConfig     = "config"
	KubeNodeDefaultDiskConfigAnnotationKey    = "node.longhorn.io/default-disks-config"
	KubeNodeDefaultNodeTagConfigAnnotationKey = "node.longhorn.io/default-node-tags"

	LastAppliedTolerationAnnotationKeySuffix = "last-applied-tolerations"

	KubernetesStatusLabel = "KubernetesStatus"
	KubernetesReplicaSet  = "ReplicaSet"
	KubernetesStatefulSet = "StatefulSet"
	RecurringJobLabel     = "RecurringJob"

	LonghornLabelKeyPrefix = "longhorn.io"

	LonghornLabelEngineImage          = "engine-image"
	LonghornLabelInstanceManager      = "instance-manager"
	LonghornLabelNode                 = "node"
	LonghornLabelDiskUUID             = "disk-uuid"
	LonghornLabelInstanceManagerType  = "instance-manager-type"
	LonghornLabelInstanceManagerImage = "instance-manager-image"
	LonghornLabelVolume               = "longhornvolume"
	LonghornLabelShareManager         = "share-manager"
	LonghornLabelShareManagerImage    = "share-manager-image"
	LonghornLabelBackingImage         = "backing-image"
	LonghornLabelBackingImageManager  = "backing-image-manager"
	LonghornLabelManagedBy            = "managed-by"
	LonghornLabelCronJobTask          = "job-task"

	LonghornLabelBackingImageDataSource = "backing-image-data-source"

	KubernetesFailureDomainRegionLabelKey = "failure-domain.beta.kubernetes.io/region"
	KubernetesFailureDomainZoneLabelKey   = "failure-domain.beta.kubernetes.io/zone"
	KubernetesTopologyRegionLabelKey      = "topology.kubernetes.io/region"
	KubernetesTopologyZoneLabelKey        = "topology.kubernetes.io/zone"

	LonghornDriverName = "driver.longhorn.io"

	DefaultDiskPrefix = "default-disk-"

	DeprecatedProvisionerName        = "rancher.io/longhorn"
	DepracatedDriverName             = "io.rancher.longhorn"
	DefaultStorageClassConfigMapName = "longhorn-storageclass"
	DefaultStorageClassName          = "longhorn"
	ControlPlaneName                 = "longhorn-manager"
)

const (
	CSIMinVersion                = "v1.14.0"
	CSIVolumeExpansionMinVersion = "v1.16.0"
	CSISnapshotterMinVersion     = "v1.17.0"

	KubernetesTopologyLabelsVersion = "v1.17.0"
)

type ReplicaMode string

const (
	ReplicaModeRW  = ReplicaMode("RW")
	ReplicaModeWO  = ReplicaMode("WO")
	ReplicaModeERR = ReplicaMode("ERR")

	EnvNodeName       = "NODE_NAME"
	EnvPodNamespace   = "POD_NAMESPACE"
	EnvPodIP          = "POD_IP"
	EnvServiceAccount = "SERVICE_ACCOUNT"

	BackupStoreTypeS3 = "s3"

	AWSIAMRoleAnnotation = "iam.amazonaws.com/role"
	AWSIAMRoleArn        = "AWS_IAM_ROLE_ARN"
	AWSAccessKey         = "AWS_ACCESS_KEY_ID"
	AWSSecretKey         = "AWS_SECRET_ACCESS_KEY"
	AWSEndPoint          = "AWS_ENDPOINTS"
	AWSCert              = "AWS_CERT"

	HTTPSProxy = "HTTPS_PROXY"
	HTTPProxy  = "HTTP_PROXY"
	NOProxy    = "NO_PROXY"

	VirtualHostedStyle = "VIRTUAL_HOSTED_STYLE"

	OptionFromBackup          = "fromBackup"
	OptionNumberOfReplicas    = "numberOfReplicas"
	OptionStaleReplicaTimeout = "staleReplicaTimeout"
	OptionBaseImage           = "baseImage"
	OptionFrontend            = "frontend"
	OptionDiskSelector        = "diskSelector"
	OptionNodeSelector        = "nodeSelector"

	// DefaultStaleReplicaTimeout in minutes. 48h by default
	DefaultStaleReplicaTimeout = "2880"

	ImageChecksumNameLength = 8
)

type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("cannot find %v", e.Name)
}

const (
	engineSuffix    = "-e"
	replicaSuffix   = "-r"
	recurringSuffix = "-c"

	// MaximumJobNameSize is calculated using
	// 1. NameMaximumLength is 40
	// 2. Recurring suffix is 2
	// 3. Maximum kubernetes name length is 63
	// 4. cronjob pod suffix is 11
	// 5. Dash and buffer for 2
	MaximumJobNameSize = 8

	engineImagePrefix          = "ei-"
	instanceManagerImagePrefix = "imi-"
	shareManagerImagePrefix    = "smi-"

	shareManagerPrefix    = "share-manager-"
	instanceManagerPrefix = "instance-manager-"
	engineManagerPrefix   = instanceManagerPrefix + "e-"
	replicaManagerPrefix  = instanceManagerPrefix + "r-"
)

func GenerateEngineNameForVolume(vName string) string {
	return vName + engineSuffix + "-" + util.RandomID()
}

func GenerateReplicaNameForVolume(vName string) string {
	return vName + replicaSuffix + "-" + util.RandomID()
}

func GetCronJobNameForVolumeAndJob(vName, job string) string {
	return vName + "-" + job + recurringSuffix
}

func GetAPIServerAddressFromIP(ip string) string {
	return net.JoinHostPort(ip, strconv.Itoa(DefaultAPIPort))
}

func GetDefaultManagerURL() string {
	return "http://longhorn-backend:" + strconv.Itoa(DefaultAPIPort) + "/v1"
}

func GetImageCanonicalName(image string) string {
	return strings.Replace(strings.Replace(image, ":", "-", -1), "/", "-", -1)
}

func GetEngineBinaryDirectoryOnHostForImage(image string) string {
	cname := GetImageCanonicalName(image)
	return filepath.Join(EngineBinaryDirectoryOnHost, cname)
}

func GetEngineBinaryDirectoryForEngineManagerContainer(image string) string {
	cname := GetImageCanonicalName(image)
	return filepath.Join(EngineBinaryDirectoryInContainer, cname)
}

func GetEngineBinaryDirectoryForReplicaManagerContainer(image string) string {
	cname := GetImageCanonicalName(image)
	return filepath.Join(filepath.Join(ReplicaHostPrefix, EngineBinaryDirectoryOnHost), cname)
}

func EngineBinaryExistOnHostForImage(image string) bool {
	st, err := os.Stat(filepath.Join(GetEngineBinaryDirectoryOnHostForImage(image), "longhorn"))
	return err == nil && !st.IsDir()
}

func GetBackingImageManagerName(image, diskUUID string) string {
	return fmt.Sprintf("backing-image-manager-%s-%s", util.GetStringChecksum(image)[:4], diskUUID[:4])
}

func GetBackingImageDirectoryName(backingImageName, backingImageUUID string) string {
	return fmt.Sprintf("%s-%s", backingImageName, backingImageUUID)
}

func GetBackingImageManagerDirectoryOnHost(diskPath string) string {
	return filepath.Join(diskPath, BackingImagesManagerDirectory)
}

func GetBackingImageDirectoryOnHost(diskPath, backingImageName, backingImageUUID string) string {
	return filepath.Join(GetBackingImageManagerDirectoryOnHost(diskPath), GetBackingImageDirectoryName(backingImageName, backingImageUUID))
}

func GetBackingImagePathForReplicaManagerContainer(diskPath, backingImageName, backingImageUUID string) string {
	return filepath.Join(ReplicaHostPrefix, GetBackingImageDirectoryOnHost(diskPath, backingImageName, backingImageUUID), BackingImageFileName)
}

var (
	LonghornSystemKey = "longhorn"
)

func GetLonghornLabelKey(name string) string {
	return fmt.Sprintf("%s/%s", LonghornLabelKeyPrefix, name)
}

func GetBaseLabelsForSystemManagedComponent() map[string]string {
	return map[string]string{GetLonghornLabelKey(LonghornLabelManagedBy): ControlPlaneName}
}

func GetLonghornLabelComponentKey() string {
	return GetLonghornLabelKey("component")
}

func GetEngineImageLabels(engineImageName string) map[string]string {
	labels := GetBaseLabelsForSystemManagedComponent()
	labels[GetLonghornLabelComponentKey()] = LonghornLabelEngineImage
	labels[GetLonghornLabelKey(LonghornLabelEngineImage)] = engineImageName
	return labels
}

// GetEIDaemonSetLabelSelector returns labels for engine image daemonset's Spec.Selector.MatchLabels
func GetEIDaemonSetLabelSelector(engineImageName string) map[string]string {
	labels := make(map[string]string)
	labels[GetLonghornLabelComponentKey()] = LonghornLabelEngineImage
	labels[GetLonghornLabelKey(LonghornLabelEngineImage)] = engineImageName
	return labels
}

func GetEngineImageComponentLabel() map[string]string {
	return map[string]string{
		GetLonghornLabelComponentKey(): LonghornLabelEngineImage,
	}
}

func GetInstanceManagerLabels(node, instanceManagerImage string, managerType InstanceManagerType) map[string]string {
	labels := GetBaseLabelsForSystemManagedComponent()
	labels[GetLonghornLabelComponentKey()] = LonghornLabelInstanceManager
	labels[GetLonghornLabelKey(LonghornLabelInstanceManagerType)] = string(managerType)
	if node != "" {
		labels[GetLonghornLabelKey(LonghornLabelNode)] = node
	}
	if instanceManagerImage != "" {
		labels[GetLonghornLabelKey(LonghornLabelInstanceManagerImage)] = GetInstanceManagerImageChecksumName(GetImageCanonicalName(instanceManagerImage))
	}

	return labels
}

func GetInstanceManagerComponentLabel() map[string]string {
	return map[string]string{
		GetLonghornLabelComponentKey(): LonghornLabelInstanceManager,
	}
}

func GetShareManagerComponentLabel() map[string]string {
	return map[string]string{
		GetLonghornLabelComponentKey(): LonghornLabelShareManager,
	}
}

func GetShareManagerInstanceLabel(name string) map[string]string {
	labels := GetBaseLabelsForSystemManagedComponent()
	labels[GetLonghornLabelKey(LonghornLabelShareManager)] = name
	return labels
}

func GetShareManagerLabels(name, image string) map[string]string {
	labels := GetBaseLabelsForSystemManagedComponent()
	labels[GetLonghornLabelComponentKey()] = LonghornLabelShareManager

	if name != "" {
		labels[GetLonghornLabelKey(LonghornLabelShareManager)] = name
	}

	if image != "" {
		labels[GetLonghornLabelKey(LonghornLabelShareManagerImage)] = GetShareManagerImageChecksumName(GetImageCanonicalName(image))
	}

	return labels
}

func GetCronJobLabels(volumeName string, job *RecurringJob) map[string]string {
	labels := GetBaseLabelsForSystemManagedComponent()
	labels[LonghornLabelVolume] = volumeName
	labels[GetLonghornLabelKey(LonghornLabelCronJobTask)] = string(job.Task)
	return labels
}

func GetCronJobPodLabels(volumeName string, job *RecurringJob) map[string]string {
	labels := make(map[string]string)
	labels[LonghornLabelVolume] = volumeName
	labels[GetLonghornLabelKey(LonghornLabelCronJobTask)] = string(job.Task)
	return labels
}

func GetBackingImageLabels() map[string]string {
	labels := GetBaseLabelsForSystemManagedComponent()
	labels[GetLonghornLabelComponentKey()] = LonghornLabelBackingImage
	return labels
}

func GetBackingImageManagerLabels(nodeID, diskUUID string) map[string]string {
	labels := GetBaseLabelsForSystemManagedComponent()
	labels[GetLonghornLabelComponentKey()] = LonghornLabelBackingImageManager
	if diskUUID != "" {
		labels[GetLonghornLabelKey(LonghornLabelDiskUUID)] = diskUUID
	}
	if nodeID != "" {
		labels[GetLonghornLabelKey(LonghornLabelNode)] = nodeID
	}
	return labels
}

func GetBackingImageDataSourceLabels(name, nodeID, diskUUID string) map[string]string {
	labels := GetBaseLabelsForSystemManagedComponent()
	labels[GetLonghornLabelComponentKey()] = LonghornLabelBackingImageDataSource
	if name != "" {
		labels[GetLonghornLabelKey(LonghornLabelBackingImageDataSource)] = name
	}
	if diskUUID != "" {
		labels[GetLonghornLabelKey(LonghornLabelDiskUUID)] = diskUUID
	}
	if nodeID != "" {
		labels[GetLonghornLabelKey(LonghornLabelNode)] = nodeID
	}
	return labels
}

func GetVolumeLabels(volumeName string) map[string]string {
	return map[string]string{
		LonghornLabelVolume: volumeName,
	}
}

func GetRegionAndZone(labels map[string]string, isUsingTopologyLabels bool) (string, string) {
	region := ""
	zone := ""
	if isUsingTopologyLabels {
		if v, ok := labels[KubernetesTopologyRegionLabelKey]; ok {
			region = v
		}
		if v, ok := labels[KubernetesTopologyZoneLabelKey]; ok {
			zone = v
		}
	} else {
		if v, ok := labels[KubernetesFailureDomainRegionLabelKey]; ok {
			region = v
		}
		if v, ok := labels[KubernetesFailureDomainZoneLabelKey]; ok {
			zone = v
		}
	}
	return region, zone
}

func GetEngineImageChecksumName(image string) string {
	return engineImagePrefix + util.GetStringChecksum(strings.TrimSpace(image))[:ImageChecksumNameLength]
}

func GetInstanceManagerImageChecksumName(image string) string {
	return instanceManagerImagePrefix + util.GetStringChecksum(strings.TrimSpace(image))[:ImageChecksumNameLength]
}

func GetShareManagerImageChecksumName(image string) string {
	return shareManagerImagePrefix + util.GetStringChecksum(strings.TrimSpace(image))[:ImageChecksumNameLength]
}

func GetShareManagerPodNameFromShareManagerName(smName string) string {
	return LonghornLabelShareManager + "-" + smName
}

func GetShareManagerNameFromShareManagerPodName(podName string) string {
	return strings.TrimPrefix(podName, LonghornLabelShareManager+"-")
}

func ValidateEngineImageChecksumName(name string) bool {
	matched, _ := regexp.MatchString(fmt.Sprintf("^%s[a-fA-F0-9]{%d}$", engineImagePrefix, ImageChecksumNameLength), name)
	return matched
}

func GetInstanceManagerName(imType InstanceManagerType) (string, error) {
	switch imType {
	case InstanceManagerTypeEngine:
		return engineManagerPrefix + util.RandomID(), nil
	case InstanceManagerTypeReplica:
		return replicaManagerPrefix + util.RandomID(), nil
	}
	return "", fmt.Errorf("cannot generate name for unknown instance manager type %v", imType)
}

func GetInstanceManagerPrefix(imType InstanceManagerType) string {
	switch imType {
	case InstanceManagerTypeEngine:
		return engineManagerPrefix
	case InstanceManagerTypeReplica:
		return replicaManagerPrefix
	}
	return ""
}

func GetReplicaDataPath(diskPath, dataDirectoryName string) string {
	return filepath.Join(diskPath, "replicas", dataDirectoryName)
}

func GetReplicaMountedDataPath(dataPath string) string {
	if !strings.HasPrefix(dataPath, ReplicaHostPrefix) {
		return filepath.Join(ReplicaHostPrefix, dataPath)
	}
	return dataPath
}

func ErrorIsNotFound(err error) bool {
	return strings.Contains(err.Error(), "cannot find")
}

func ErrorAlreadyExists(err error) bool {
	return strings.Contains(err.Error(), "already exists")
}

func ValidateReplicaCount(count int) error {
	if count < 1 || count > 20 {
		return fmt.Errorf("replica count value must between 1 to 20")
	}
	return nil
}

func ValidateReplicaAutoBalance(option ReplicaAutoBalance) error {
	switch option {
	case ReplicaAutoBalanceIgnored,
		ReplicaAutoBalanceDisabled,
		ReplicaAutoBalanceLeastEffort,
		ReplicaAutoBalanceBestEffort:
		return nil
	default:
		return fmt.Errorf("invalid replica auto-balance option: %v", option)
	}
}

func ValidateDataLocality(mode DataLocality) error {
	if mode != DataLocalityDisabled && mode != DataLocalityBestEffort {
		return fmt.Errorf("invalid data locality mode: %v", mode)
	}
	return nil
}

func ValidateAccessMode(mode AccessMode) error {
	if mode != AccessModeReadWriteMany && mode != AccessModeReadWriteOnce {
		return fmt.Errorf("invalid access mode: %v", mode)
	}
	return nil
}

func GetDaemonSetNameFromEngineImageName(engineImageName string) string {
	return "engine-image-" + engineImageName
}

func GetEngineImageNameFromDaemonSetName(dsName string) string {
	return strings.TrimPrefix(dsName, "engine-image-")
}

func LabelsToString(labels map[string]string) string {
	res := ""
	for k, v := range labels {
		res += fmt.Sprintf("%s=%s,", k, v)
	}
	res = strings.TrimSuffix(res, ",")
	return res
}

func CreateDisksFromAnnotation(annotation string) (map[string]DiskSpec, error) {
	validDisks := map[string]DiskSpec{}
	existFsid := map[string]string{}

	disks, err := UnmarshalToDisks(annotation)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal the default disks annotation")
	}
	for _, disk := range disks {
		if disk.Path == "" {
			return nil, fmt.Errorf("invalid disk %+v", disk)
		}
		diskInfo, err := util.GetDiskInfo(disk.Path)
		if err != nil {
			return nil, err
		}
		for _, vDisk := range validDisks {
			if vDisk.Path == disk.Path {
				return nil, fmt.Errorf("duplicate disk path %v", disk.Path)
			}
		}

		// Set to default disk name
		if disk.Name == "" {
			disk.Name = DefaultDiskPrefix + diskInfo.Fsid
		}

		if _, exist := existFsid[diskInfo.Fsid]; exist {
			return nil, fmt.Errorf(
				"the disk %v is the same"+
					"file system with %v, fsid %v",
				disk.Path, existFsid[diskInfo.Fsid],
				diskInfo.Fsid)
		}

		existFsid[diskInfo.Fsid] = disk.Path

		if disk.StorageReserved < 0 || disk.StorageReserved > diskInfo.StorageMaximum {
			return nil, fmt.Errorf("the storageReserved setting of disk %v is not valid, should be positive and no more than storageMaximum and storageAvailable", disk.Path)
		}
		tags, err := util.ValidateTags(disk.Tags)
		if err != nil {
			return nil, err
		}
		disk.Tags = tags
		_, exists := validDisks[disk.Name]
		if exists {
			return nil, fmt.Errorf("the disk name %v has duplicated", disk.Name)
		}
		validDisks[disk.Name] = disk.DiskSpec
	}

	return validDisks, nil
}

func GetNodeTagsFromAnnotation(annotation string) ([]string, error) {
	nodeTags, err := UnmarshalToNodeTags(annotation)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal the node tag annotation")
	}
	validNodeTags, err := util.ValidateTags(nodeTags)
	if err != nil {
		return nil, err
	}

	return validNodeTags, nil
}

type DiskSpecWithName struct {
	DiskSpec
	Name string `json:"name"`
}

// UnmarshalToDisks input format should be:
// `[{"path":"/mnt/disk1","allowScheduling":false},
//   {"path":"/mnt/disk2","allowScheduling":false,"storageReserved":1024,"tags":["ssd","fast"]}]`
func UnmarshalToDisks(s string) (ret []DiskSpecWithName, err error) {
	if err := json.Unmarshal([]byte(s), &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// UnmarshalToNodeTags input format should be:
// `["worker1","enabled"]`
func UnmarshalToNodeTags(s string) ([]string, error) {
	var res []string
	if err := json.Unmarshal([]byte(s), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func CreateDefaultDisk(dataPath string) (map[string]DiskSpec, error) {
	if err := util.CreateDiskPathReplicaSubdirectory(dataPath); err != nil {
		return nil, err
	}
	diskInfo, err := util.GetDiskInfo(dataPath)
	if err != nil {
		return nil, err
	}
	return map[string]DiskSpec{
		DefaultDiskPrefix + diskInfo.Fsid: {
			Path:              diskInfo.Path,
			AllowScheduling:   true,
			EvictionRequested: false,
			StorageReserved:   diskInfo.StorageMaximum * 30 / 100,
		},
	}, nil
}

func ValidateCPUReservationValues(engineManagerCPUStr, replicaManagerCPUStr string) error {
	engineManagerCPU, err := strconv.Atoi(engineManagerCPUStr)
	if err != nil {
		return fmt.Errorf("guaranteed/requested engine manager CPU value %v is not int: %v", engineManagerCPUStr, err)
	}
	replicaManagerCPU, err := strconv.Atoi(replicaManagerCPUStr)
	if err != nil {
		return fmt.Errorf("guaranteed/requested replica manager CPU value %v is not int: %v", replicaManagerCPUStr, err)
	}
	if engineManagerCPU+replicaManagerCPU < 0 || engineManagerCPU+replicaManagerCPU > 40 {
		return fmt.Errorf("the requested engine manager CPU and replica manager CPU are %v%% and %v%% of a node total CPU, respectively. The sum should not be smaller than 0%% or greater than 40%%", engineManagerCPU, replicaManagerCPU)
	}
	return nil
}
