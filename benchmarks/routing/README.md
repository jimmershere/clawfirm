# Routing benchmark

Held-out test set of `(message, expected_class, expected_sensitivity, expected_complexity)` tuples. Run the policy-judge against each, score F1 per rubric, and produce a confusion matrix.

The test set is versioned with the rubrics in `services/policy-judge/policy_judge/prompts/`. When a rubric changes, the test set must be updated and re-scored before the change ships.
