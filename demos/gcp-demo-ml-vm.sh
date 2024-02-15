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

# VPC

# dataset
gcloud compute networks create gcp-machine-learning-dataset --project=gcp-ibngctopoc-nprd-72084 --description=AWI\ demo\ network --subnet-mode=custom --mtu=1460 --bgp-routing-mode=regional

gcloud compute networks subnets create gcp-machine-learning-dataset-subnet1 --project=gcp-ibngctopoc-nprd-72084 --range=10.70.0.0/16 --stack-type=IPV4_ONLY --network=gcp-machine-learning-dataset --region=us-west1

gcloud compute firewall-rules create gcp-machine-learning-dataset-allow-ssh --project=gcp-ibngctopoc-nprd-72084 --network=projects/gcp-ibngctopoc-nprd-72084/global/networks/gcp-machine-learning-dataset --description=Allows\ TCP\ connections\ from\ any\ source\ to\ any\ instance\ on\ the\ network\ using\ port\ 22. --direction=INGRESS --priority=65534 --source-ranges=0.0.0.0/0 --action=ALLOW --rules=tcp:22


gcloud compute instances create machine-learning-dataset \
    --project=gcp-ibngctopoc-nprd-72084 \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --network-interface=network-tier=STANDARD,subnet=gcp-machine-learning-dataset-subnet1 \
    --maintenance-policy=MIGRATE \
    --provisioning-model=STANDARD \
    --service-account=475295088090-compute@developer.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append \
    --create-disk=auto-delete=yes,boot=yes,device-name=development-database-2,image=projects/debian-cloud/global/images/debian-11-bullseye-v20230306,mode=rw,size=10,type=projects/gcp-ibngctopoc-nprd-72084/zones/us-west1-b/diskTypes/pd-balanced \
    --no-shielded-secure-boot \
    --shielded-vtpm \
    --shielded-integrity-monitoring \
    --labels=environment=dataset,ec-src=vm_add-gcloud \
    --reservation-affinity=any

# training
gcloud compute networks create gcp-machine-learning-training --project=gcp-ibngctopoc-nprd-72084 --description=AWI\ demo\ network --subnet-mode=custom --mtu=1460 --bgp-routing-mode=regional

gcloud compute networks subnets create gcp-machine-learning-training-subnet1 --project=gcp-ibngctopoc-nprd-72084 --range=10.80.0.0/16 --stack-type=IPV4_ONLY --network=gcp-machine-learning-training --region=us-west1

gcloud compute networks subnets create gcp-machine-learning-training-subnet-fe --project=gcp-ibngctopoc-nprd-72084 --range=10.81.0.0/16 --stack-type=IPV4_ONLY --network=gcp-machine-learning-training --region=us-west1

gcloud compute firewall-rules create gcp-machine-learning-training-allow-ssh --project=gcp-ibngctopoc-nprd-72084 --network=projects/gcp-ibngctopoc-nprd-72084/global/networks/gcp-machine-learning-training --description=Allows\ TCP\ connections\ from\ any\ source\ to\ any\ instance\ on\ the\ network\ using\ port\ 22. --direction=INGRESS --priority=65534 --source-ranges=0.0.0.0/0 --action=ALLOW --rules=tcp:22


gcloud compute instances create machine-learning-training-1 \
    --project=gcp-ibngctopoc-nprd-72084 \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --network-interface=network-tier=STANDARD,subnet=gcp-machine-learning-training-subnet1 \
    --maintenance-policy=MIGRATE \
    --provisioning-model=STANDARD \
    --service-account=475295088090-compute@developer.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append \
    --create-disk=auto-delete=yes,boot=yes,device-name=development-database-2,image=projects/debian-cloud/global/images/debian-11-bullseye-v20230306,mode=rw,size=10,type=projects/gcp-ibngctopoc-nprd-72084/zones/us-west1-b/diskTypes/pd-balanced \
    --no-shielded-secure-boot \
    --shielded-vtpm \
    --shielded-integrity-monitoring \
    --labels=environment=training,app_type=backend,ec-src=vm_add-gcloud \
    --reservation-affinity=any

gcloud compute instances create machine-learning-training-2 \
    --project=gcp-ibngctopoc-nprd-72084 \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --network-interface=network-tier=STANDARD,subnet=gcp-machine-learning-training-subnet1 \
    --maintenance-policy=MIGRATE \
    --provisioning-model=STANDARD \
    --service-account=475295088090-compute@developer.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append \
    --create-disk=auto-delete=yes,boot=yes,device-name=development-database-2,image=projects/debian-cloud/global/images/debian-11-bullseye-v20230306,mode=rw,size=10,type=projects/gcp-ibngctopoc-nprd-72084/zones/us-west1-b/diskTypes/pd-balanced \
    --no-shielded-secure-boot \
    --shielded-vtpm \
    --shielded-integrity-monitoring \
    --labels=environment=training,app_type=backend,ec-src=vm_add-gcloud \
    --reservation-affinity=any

gcloud compute instances create machine-learning-frontend \
    --project=gcp-ibngctopoc-nprd-72084 \
    --zone=us-west1-b \
    --machine-type=e2-micro \
    --network-interface=network-tier=STANDARD,subnet=gcp-machine-learning-training-subnet-fe \
    --maintenance-policy=MIGRATE \
    --provisioning-model=STANDARD \
    --service-account=475295088090-compute@developer.gserviceaccount.com \
    --scopes=https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring.write,https://www.googleapis.com/auth/servicecontrol,https://www.googleapis.com/auth/service.management.readonly,https://www.googleapis.com/auth/trace.append \
    --create-disk=auto-delete=yes,boot=yes,device-name=development-database-2,image=projects/debian-cloud/global/images/debian-11-bullseye-v20230306,mode=rw,size=10,type=projects/gcp-ibngctopoc-nprd-72084/zones/us-west1-b/diskTypes/pd-balanced \
    --no-shielded-secure-boot \
    --shielded-vtpm \
    --shielded-integrity-monitoring \
    --labels=environment=training,app_type=frontend,ec-src=vm_add-gcloud \
    --reservation-affinity=any