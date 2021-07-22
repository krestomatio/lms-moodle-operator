.PHONY: changelog
changelog: jx-changelog ## Generate changelog
	@echo "+ $@"

.PHONY: release
release: skopeo-copy ## Run release tasks
	@echo "+ $@"

.PHONY: promote
promote: git ## Promote release
	@echo "+ $@"
