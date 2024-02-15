# Copyright (c) 2023 Cisco Systems, Inc. and its affiliates
# All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http:www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0


# TODO add VPC creation here as well

echo "Deploying development-dashboard-1 instance"
gcloud compute instances create development-dashboard-1 \
    --project=gcp-ibngctopoc-nprd-72084 \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --network-interface=network-tier=STANDARD,subnet=development-subnet-1 \
    --maintenance-policy=MIGRATE \
    --provisioning-model=STANDARD \
    --service-account=475295088090-compute@developer.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append \
    --create-disk=auto-delete=yes,boot=yes,device-name=development-dashboard-1,image=projects/debian-cloud/global/images/debian-11-bullseye-v20230306,mode=rw,size=10,type=projects/gcp-ibngctopoc-nprd-72084/zones/us-west1-b/diskTypes/pd-balanced \
    --no-shielded-secure-boot \
    --shielded-vtpm \
    --shielded-integrity-monitoring \
    --labels=env=development,app_type=dashboard,ec-src=vm_add-gcloud \
    --reservation-affinity=any

echo "Deploying development-database-1 instance"
gcloud compute instances create development-database-1 \
    --project=gcp-ibngctopoc-nprd-72084 \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --network-interface=network-tier=STANDARD,subnet=development-subnet-1 \
    --maintenance-policy=MIGRATE \
    --provisioning-model=STANDARD \
    --service-account=475295088090-compute@developer.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append \
    --create-disk=auto-delete=yes,boot=yes,device-name=development-database-1,image=projects/debian-cloud/global/images/debian-11-bullseye-v20230306,mode=rw,size=10,type=projects/gcp-ibngctopoc-nprd-72084/zones/us-west1-b/diskTypes/pd-balanced \
    --no-shielded-secure-boot \
    --shielded-vtpm \
    --shielded-integrity-monitoring \
    --labels=env=development,app_type=database,ec-src=vm_add-gcloud \
    --reservation-affinity=any

echo "Deploying development-database-2 instance"
gcloud compute instances create development-database-2 \
    --project=gcp-ibngctopoc-nprd-72084 \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --network-interface=network-tier=STANDARD,subnet=development-subnet-1 \
    --maintenance-policy=MIGRATE \
    --provisioning-model=STANDARD \
    --service-account=475295088090-compute@developer.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append \
    --create-disk=auto-delete=yes,boot=yes,device-name=development-database-2,image=projects/debian-cloud/global/images/debian-11-bullseye-v20230306,mode=rw,size=10,type=projects/gcp-ibngctopoc-nprd-72084/zones/us-west1-b/diskTypes/pd-balanced \
    --no-shielded-secure-boot \
    --shielded-vtpm \
    --shielded-integrity-monitoring \
    --labels=env=development,app_type=database,ec-src=vm_add-gcloud \
    --reservation-affinity=any


echo "Deploying staging-dashboard-1 instance"
gcloud compute instances create staging-dashboard-1 \
    --project=gcp-ibngctopoc-nprd-72084 \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --network-interface=network-tier=STANDARD,subnet=staging-subnet-1 \
    --maintenance-policy=MIGRATE \
    --provisioning-model=STANDARD \
    --service-account=475295088090-compute@developer.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append \
    --create-disk=auto-delete=yes,boot=yes,device-name=staging-dashboard-1,image=projects/debian-cloud/global/images/debian-11-bullseye-v20230306,mode=rw,size=10,type=projects/gcp-ibngctopoc-nprd-72084/zones/us-west1-b/diskTypes/pd-balanced \
    --no-shielded-secure-boot \
    --shielded-vtpm \
    --shielded-integrity-monitoring \
    --labels=env=staging,app_type=dashboard,ec-src=vm_add-gcloud \
    --reservation-affinity=any

echo "Deploying staging-database-1 instance"
gcloud compute instances create staging-database-1 \
    --project=gcp-ibngctopoc-nprd-72084 \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --network-interface=network-tier=STANDARD,subnet=staging-subnet-1 \
    --maintenance-policy=MIGRATE \
    --provisioning-model=STANDARD \
    --service-account=475295088090-compute@developer.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append \
    --create-disk=auto-delete=yes,boot=yes,device-name=staging-database-1,image=projects/debian-cloud/global/images/debian-11-bullseye-v20230306,mode=rw,size=10,type=projects/gcp-ibngctopoc-nprd-72084/zones/us-west1-b/diskTypes/pd-balanced \
    --no-shielded-secure-boot \
    --shielded-vtpm \
    --shielded-integrity-monitoring \
    --labels=env=staging,app_type=database,ec-src=vm_add-gcloud \
    --reservation-affinity=any

echo "Deploying staging-database-2 instance"
gcloud compute instances create staging-database-2 \
    --project=gcp-ibngctopoc-nprd-72084 \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --network-interface=network-tier=STANDARD,subnet=staging-subnet-1 \
    --maintenance-policy=MIGRATE \
    --provisioning-model=STANDARD \
    --service-account=475295088090-compute@developer.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append \
    --create-disk=auto-delete=yes,boot=yes,device-name=staging-database-2,image=projects/debian-cloud/global/images/debian-11-bullseye-v20230306,mode=rw,size=10,type=projects/gcp-ibngctopoc-nprd-72084/zones/us-west1-b/diskTypes/pd-balanced \
    --no-shielded-secure-boot \
    --shielded-vtpm \
    --shielded-integrity-monitoring \
    --labels=env=staging,app_type=database,ec-src=vm_add-gcloud \
    --reservation-affinity=any