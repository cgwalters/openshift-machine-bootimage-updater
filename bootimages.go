package main

import (
	"strings"
	"encoding/json"
)

type rhcosBootimage struct {
	AMIs map[string]struct {
		HVM string `json:"hvm"`
	} `json:"amis"`
	Azure struct {
		Image string `json:"image"`
		URL   string `json:"url"`
	}
	GCP struct {
		Image   string `json:"image"`
		Project string `json:"project"`
		URL     string `json:"url"`
	}
	BaseURI string `json:"baseURI"`
	Images  struct {
		QEMU struct {
			Path               string `json:"path"`
			SHA256             string `json:"sha256"`
			UncompressedSHA256 string `json:"uncompressed-sha256"`
		} `json:"qemu"`
		OpenStack struct {
			Path               string `json:"path"`
			SHA256             string `json:"sha256"`
			UncompressedSHA256 string `json:"uncompressed-sha256"`
		} `json:"openstack"`
		VMware struct {
			Path   string `json:"path"`
			SHA256 string `json:"sha256"`
		} `json:"vmware"`
	} `json:"images"`
	OSTreeVersion string `json:"ostree-version"`
}

// https://github.com/openshift/installer/blob/f5ba6239853f0904704c04d8b1c04c78172f1141/data/data/rhcos.json
var rhcos46 = `
{
    "amis": {
        "af-south-1": {
            "hvm": "ami-09921c9c1c36e695c"
        },
        "ap-east-1": {
            "hvm": "ami-01ee8446e9af6b197"
        },
        "ap-northeast-1": {
            "hvm": "ami-04e5b5722a55846ea"
        },
        "ap-northeast-2": {
            "hvm": "ami-0fdc25c8a0273a742"
        },
        "ap-south-1": {
            "hvm": "ami-09e3deb397cc526a8"
        },
        "ap-southeast-1": {
            "hvm": "ami-0630e03f75e02eec4"
        },
        "ap-southeast-2": {
            "hvm": "ami-069450613262ba03c"
        },
        "ca-central-1": {
            "hvm": "ami-012518cdbd3057dfd"
        },
        "eu-central-1": {
            "hvm": "ami-0bd7175ff5b1aef0c"
        },
        "eu-north-1": {
            "hvm": "ami-06c9ec42d0a839ad2"
        },
        "eu-south-1": {
            "hvm": "ami-0614d7440a0363d71"
        },
        "eu-west-1": {
            "hvm": "ami-01b89df58b5d4d5fa"
        },
        "eu-west-2": {
            "hvm": "ami-06f6e31ddd554f89d"
        },
        "eu-west-3": {
            "hvm": "ami-0dc82e2517ded15a1"
        },
        "me-south-1": {
            "hvm": "ami-07d181e3aa0f76067"
        },
        "sa-east-1": {
            "hvm": "ami-0cd44e6dd20e6c7fa"
        },
        "us-east-1": {
            "hvm": "ami-04a16d506e5b0e246"
        },
        "us-east-2": {
            "hvm": "ami-0a1f868ad58ea59a7"
        },
        "us-west-1": {
            "hvm": "ami-0a65d76e3a6f6622f"
        },
        "us-west-2": {
            "hvm": "ami-0dd9008abadc519f1"
        }
    },
    "azure": {
        "image": "rhcos-46.82.202011260640-0-azure.x86_64.vhd",
        "url": "https://rhcos.blob.core.windows.net/imagebucket/rhcos-46.82.202011260640-0-azure.x86_64.vhd"
    },
    "baseURI": "https://releases-art-rhcos.svc.ci.openshift.org/art/storage/releases/rhcos-4.6/46.82.202011260640-0/x86_64/",
    "buildid": "46.82.202011260640-0",
    "gcp": {
        "image": "rhcos-46-82-202011260640-0-gcp-x86-64",
        "project": "rhcos-cloud",
        "url": "https://storage.googleapis.com/rhcos/rhcos/rhcos-46-82-202011260640-0-gcp-x86-64.tar.gz"
    },
    "images": {
        "aws": {
            "path": "rhcos-46.82.202011260640-0-aws.x86_64.vmdk.gz",
            "sha256": "6e0f720077ac20fae46c16159d03fd51f66cd318155dcd14f0f5d7dc79bfd8ad",
            "size": 900764628,
            "uncompressed-sha256": "b38713b80bbcc41d18efcec93d241da48533225e0ce073c2a26235993ffc4166",
            "uncompressed-size": 919795200
        },
        "azure": {
            "path": "rhcos-46.82.202011260640-0-azure.x86_64.vhd.gz",
            "sha256": "5255c675ecff4f8932db24b68a7ef451dfc4e502326128931ccd39acf27c6c77",
            "size": 901953565,
            "uncompressed-sha256": "0fd7ab096c7ac4b8d807371d1c4a13b3a665b5d2938ff814bd3e46257193941b",
            "uncompressed-size": 17179869696
        },
        "gcp": {
            "path": "rhcos-46.82.202011260640-0-gcp.x86_64.tar.gz",
            "sha256": "f470a98c1fc7a4e0226ed4258fdd0b6edcddbb2356e7992bea4b98843cd94d9a",
            "size": 887140241
        },
        "live-initramfs": {
            "path": "rhcos-46.82.202011260640-0-live-initramfs.x86_64.img",
            "sha256": "8ff220c6f4bbca35dd071c5f1802566290b70081faeb06fb07f72c4583f8facf"
        },
        "live-iso": {
            "path": "rhcos-46.82.202011260640-0-live.x86_64.iso",
            "sha256": "161046d6275cda89a3f44582bc2b850cd09b987807d3b20090eeab279b7e0550"
        },
        "live-kernel": {
            "path": "rhcos-46.82.202011260640-0-live-kernel-x86_64",
            "sha256": "9bbb496410c04f54d30eb56f76f769bccec7e81b91271ae72261ecc54db9c63d"
        },
        "live-rootfs": {
            "path": "rhcos-46.82.202011260640-0-live-rootfs.x86_64.img",
            "sha256": "007f9c376e1282b77ed662509a3a81b624dfcf28aa971a9e80b33fc06cf9c4fd"
        },
        "metal": {
            "path": "rhcos-46.82.202011260640-0-metal.x86_64.raw.gz",
            "sha256": "bcd9871c373c02898675ec407e83f9ccd107077ff30008ecb601c9869c4a501a",
            "size": 888607273,
            "uncompressed-sha256": "a769991d850064ef078ba4dae9050e34a820bf9eba792a2b7922fb7a15315b57",
            "uncompressed-size": 3555721216
        },
        "metal4k": {
            "path": "rhcos-46.82.202011260640-0-metal4k.x86_64.raw.gz",
            "sha256": "ff4ed9d58a286185abdc4be0faaa4cd621c088121133b9b5952d2da320149fad",
            "size": 886276538,
            "uncompressed-sha256": "1e57d643969306e84b945c6355eed2793daa3e7675b0daa63defd6f5f8b5a43d",
            "uncompressed-size": 3555721216
        },
        "openstack": {
            "path": "rhcos-46.82.202011260640-0-openstack.x86_64.qcow2.gz",
            "sha256": "a8a28cfe5f5e5dadedb3442afcb447f85bddf2e82dcd558813a985a4d495782a",
            "size": 887473972,
            "uncompressed-sha256": "2bd648e09f086973accd8ac1e355ce0fcd7dfcc16bc9708c938801fcf10e219e",
            "uncompressed-size": 2254045184
        },
        "ostree": {
            "path": "rhcos-46.82.202011260640-0-ostree.x86_64.tar",
            "sha256": "1b15017685292447ebb2dc8bfe8a4f216ea1c42f871f666c65d4b98e9998d549",
            "size": 802641920
        },
        "qemu": {
            "path": "rhcos-46.82.202011260640-0-qemu.x86_64.qcow2.gz",
            "sha256": "0ea6f0852c3e0f8e4182e705561fdccd998503ef423c441e182241cd6a278730",
            "size": 888353176,
            "uncompressed-sha256": "99928ff40c2d8e3aa358d9bd453102e3d1b5e9694fb5d54febc56e275f35da51",
            "uncompressed-size": 2290221056
        },
        "vmware": {
            "path": "rhcos-46.82.202011260640-0-vmware.x86_64.ova",
            "sha256": "d73c7bdfd22e5a1231e23663266744c98c10de6f08f6d1fc718e7f72d0490c4a",
            "size": 919808000
        }
    },
    "oscontainer": {
        "digest": "sha256:54f1ab852c6592d6ca453cabb855b5c9f5ced35b1aadbb5d8161d4ca2d3623cc",
        "image": "quay.io/openshift-release-dev/ocp-v4.0-art-dev"
    },
    "ostree-commit": "cb0327325553e6922ff25822ea7eb1a2ec213e70c7cf6880965e7e2bb5ee7dea",
    "ostree-version": "46.82.202011260640-0"
}
`

func bootimageFromChannel(channel string) *rhcosBootimage {
	var r rhcosBootimage
	if strings.HasSuffix(channel, "-4.6") {
		err := json.Unmarshal([]byte(rhcos46), &r)
		if err != nil {
			panic(err)
		}
		return &r
	}
	return nil
}
