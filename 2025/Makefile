# Create a day directory from `template/`.
# Usage: `make day 3` or `make day DAY=3`
.PHONY: day

# If the second goal is provided (e.g. `make day 3`), use it as DAY.
DAY_FROM_GOALS := $(word 2,$(MAKECMDGOALS))
ifeq ($(strip $(DAY_FROM_GOALS)),)
  DAY := $(DAY)
else
  DAY := $(DAY_FROM_GOALS)
endif

day:
	@if [ -z "$(DAY)" ]; then \
		echo "Usage: make day <N> or make day DAY=<N> (where N is 1-31)"; exit 2; \
	fi
	@python scripts/start-day.py $(DAY)

# Catch-all rule to prevent make from trying to build numeric goals (like `3`).
# This defines an empty rule for any single-word goal not otherwise defined.
%:
	@:

