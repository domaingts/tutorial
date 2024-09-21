NAME = sing-box
COMMIT = $(shell git rev-parse --short HEAD)
TAGS_GO120 = with_gvisor,with_dhcp,with_wireguard,with_reality_server,with_clash_api,with_quic,with_utls
TAGS ?= $(TAGS_GO120)
TAGS_TEST ?= with_gvisor,with_quic,with_wireguard,with_grpc,with_ech,with_utls,with_reality_server

GOHOSTOS = $(shell go env GOHOSTOS)
GOHOSTARCH = $(shell go env GOHOSTARCH)
VERSION=$(shell go run ./cmd/internal/read_tag)

PARAMS = -v -trimpath -ldflags "-X 'github.com/sagernet/sing-box/constant.Version=$(VERSION)' -s -w -buildid="
MAIN_PARAMS = $(PARAMS) -tags $(TAGS)
MAIN = ./cmd/sing-box
PREFIX ?= $(shell go env GOPATH)

.PHONY: test release docs build

build:
	go build $(MAIN_PARAMS) $(MAIN)

pack:
	tar czvf sing-box.tar.gz sing-box

ci_build_go120:
	go build $(PARAMS) $(MAIN)
	go build $(PARAMS) -tags "$(TAGS_GO120)" $(MAIN)
	go build $(PARAMS) -tags "$(TAGS_GO120)" $(MAIN)

ci_build:
	go build $(PARAMS) $(MAIN)
	go build $(MAIN_PARAMS) $(MAIN)

install:
	go build -o $(PREFIX)/bin/$(NAME) $(MAIN_PARAMS) $(MAIN)

fmt:
	@gofumpt -l -w .
	@gofmt -s -w .
	@gci write --custom-order -s standard -s "prefix(github.com/sagernet/)" -s "default" .

fmt_install:
	go install -v mvdan.cc/gofumpt@latest
	go install -v github.com/daixiang0/gci@latest

lint:
	GOOS=linux golangci-lint run ./...
	GOOS=android golangci-lint run ./...
	GOOS=windows golangci-lint run ./...
	GOOS=darwin golangci-lint run ./...
	GOOS=freebsd golangci-lint run ./...

lint_install:
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest

proto:
	@go run ./cmd/internal/protogen
	@gofumpt -l -w .
	@gofumpt -l -w .

proto_install:
	go install -v google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

release:
	go run ./cmd/internal/build goreleaser release --clean --skip publish
	mkdir dist/release
	mv dist/*.tar.gz \
		dist/*.zip \
		dist/*.deb \
		dist/*.rpm \
		dist/*_amd64.pkg.tar.zst \
		dist/*_amd64v3.pkg.tar.zst \
		dist/*_arm64.pkg.tar.zst \
		dist/release
	ghr --replace --draft --prerelease -p 3 "v${VERSION}" dist/release
	rm -r dist/release

release_repo:
	go run ./cmd/internal/build goreleaser release -f .goreleaser.fury.yaml --clean

release_install:
	go install -v github.com/tcnksm/ghr@latest

update_android_version:
	go run ./cmd/internal/update_android_version

build_android:
	cd ../sing-box-for-android && ./gradlew :app:clean :app:assemblePlayRelease :app:assembleOtherRelease && ./gradlew --stop

upload_android:
	mkdir -p dist/release_android
	cp ../sing-box-for-android/app/build/outputs/apk/play/release/*.apk dist/release_android
	cp ../sing-box-for-android/app/build/outputs/apk/other/release/*-universal.apk dist/release_android
	ghr --replace --draft --prerelease -p 3 "v${VERSION}" dist/release_android
	rm -rf dist/release_android

release_android: lib_android update_android_version build_android upload_android

publish_android:
	cd ../sing-box-for-android && ./gradlew :app:publishPlayReleaseBundle

publish_android_appcenter:
	cd ../sing-box-for-android && ./gradlew :app:appCenterAssembleAndUploadPlayRelease


# TODO: find why and remove `-destination 'generic/platform=iOS'`
build_ios:
	cd ../sing-box-for-apple && \
	rm -rf build/SFI.xcarchive && \
	xcodebuild archive -scheme SFI -configuration Release -destination 'generic/platform=iOS' -archivePath build/SFI.xcarchive -allowProvisioningUpdates

upload_ios_app_store:
	cd ../sing-box-for-apple && \
	xcodebuild -exportArchive -archivePath build/SFI.xcarchive -exportOptionsPlist SFI/Upload.plist -allowProvisioningUpdates

release_ios: build_ios upload_ios_app_store

build_macos:
	cd ../sing-box-for-apple && \
	rm -rf build/SFM.xcarchive && \
	xcodebuild archive -scheme SFM -configuration Release -archivePath build/SFM.xcarchive -allowProvisioningUpdates

upload_macos_app_store:
	cd ../sing-box-for-apple && \
	xcodebuild -exportArchive -archivePath build/SFM.xcarchive -exportOptionsPlist SFI/Upload.plist -allowProvisioningUpdates

release_macos: build_macos upload_macos_app_store

build_macos_standalone:
	cd ../sing-box-for-apple && \
	rm -rf build/SFM.System.xcarchive && \
	xcodebuild archive -scheme SFM.System -configuration Release -archivePath build/SFM.System.xcarchive -allowProvisioningUpdates

build_macos_dmg:
	rm -rf dist/SFM
	mkdir -p dist/SFM
	cd ../sing-box-for-apple && \
	rm -rf build/SFM.System && \
	rm -rf build/SFM.dmg && \
	xcodebuild -exportArchive \
		-archivePath "build/SFM.System.xcarchive" \
		-exportOptionsPlist SFM.System/Export.plist -allowProvisioningUpdates \
		-exportPath "build/SFM.System" && \
	create-dmg \
		--volname "sing-box" \
		--volicon "build/SFM.System/SFM.app/Contents/Resources/AppIcon.icns" \
		--icon "SFM.app" 0 0 \
 		--hide-extension "SFM.app" \
 		--app-drop-link 0 0 \
 		--skip-jenkins \
		--notarize "notarytool-password" \
		"../sing-box/dist/SFM/SFM.dmg" "build/SFM.System/SFM.app"

upload_macos_dmg:
	cd dist/SFM && \
	cp SFM.dmg "SFM-${VERSION}-universal.dmg" && \
	ghr --replace --draft --prerelease "v${VERSION}" "SFM-${VERSION}-universal.dmg"

release_macos_standalone: build_macos_standalone build_macos_dmg upload_macos_dmg

build_tvos:
	cd ../sing-box-for-apple && \
	rm -rf build/SFT.xcarchive && \
	xcodebuild archive -scheme SFT -configuration Release -archivePath build/SFT.xcarchive -allowProvisioningUpdates

upload_tvos_app_store:
	cd ../sing-box-for-apple && \
	xcodebuild -exportArchive -archivePath "build/SFT.xcarchive" -exportOptionsPlist SFI/Upload.plist -allowProvisioningUpdates

release_tvos: build_tvos upload_tvos_app_store

update_apple_version:
	go run ./cmd/internal/update_apple_version

release_apple: lib_ios update_apple_version release_ios release_macos release_tvos

release_apple_beta: update_apple_version release_ios release_macos release_tvos

test:
	@go test -v ./... && \
	cd test && \
	go mod tidy && \
	go test -v -tags "$(TAGS_TEST)" .

test_stdio:
	@go test -v ./... && \
	cd test && \
	go mod tidy && \
	go test -v -tags "$(TAGS_TEST),force_stdio" .

lib_android:
	go run ./cmd/internal/build_libbox -target android

lib_ios:
	go run ./cmd/internal/build_libbox -target ios

lib:
	go run ./cmd/internal/build_libbox -target android
	go run ./cmd/internal/build_libbox -target ios

lib_install:
	go install -v github.com/sagernet/gomobile/cmd/gomobile@v0.1.4
	go install -v github.com/sagernet/gomobile/cmd/gobind@v0.1.4

docs:
	venv/bin/mkdocs serve

publish_docs:
	venv/bin/mkdocs gh-deploy -m "Update" --force --ignore-version --no-history

docs_install:
	python -m venv venv
	source ./venv/bin/activate && pip install --force-reinstall mkdocs-material=="9.*" mkdocs-static-i18n=="1.2.*"

clean:
	rm -rf bin dist sing-box
	rm -f $(shell go env GOPATH)/sing-box

update:
	git fetch
	git reset FETCH_HEAD --hard
	git clean -fdx