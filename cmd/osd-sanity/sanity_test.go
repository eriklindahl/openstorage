/*
Copyright 2017 Portworx

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package sanity

import (
	"flag"
	"fmt"
	"testing"

	"github.com/libopenstorage/openstorage/pkg/sanity"
)

const (
	prefix string = "osd."
)

var (
	VERSION               = "(dev)"
	endpoint              string
	driver                string
	cloudBackupConfigPath string
	version               bool
)

func init() {
	flag.StringVar(&endpoint, prefix+"endpoint", "", "OSD endpoint")
	flag.StringVar(&driver, prefix+"driver", "", "Volume Driver")
	flag.StringVar(&cloudBackupConfigPath, prefix+"cloudbackupconfig", "", "Test Cloud Backup/Restore config file , optional")
	flag.BoolVar(&version, prefix+"version", false, "Version of this program")
	flag.Parse()
}

func TestSanity(t *testing.T) {
	if version {
		fmt.Printf("Version = %s\n", VERSION)
		return
	}
	if len(endpoint) == 0 {
		t.Fatalf("--%s.endpoint must be provided with an OSD endpoint", prefix)
	}
	if len(driver) == 0 {
		t.Fatalf("--%s.driver must be provided with a Volume driver", prefix)
	}
	if len(cloudBackupConfigPath) == 0 {
		t.Logf("No Cloud backup config file provided , Cloud Backup & Restore Tests will be skipped")
	}

	sanity.Test(t, endpoint, driver, cloudBackupConfigPath)
}
