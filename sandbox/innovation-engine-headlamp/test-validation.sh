#!/bin/bash
echo "Testing with no credentials..."
rm -f .env
npm run validate-env:prod
prod_no_creds_status=$?
echo "Production mode with no credentials exit code: $prod_no_creds_status (should be non-zero)"

npm run validate-env:dev
dev_no_creds_status=$?
echo "Development mode with no credentials exit code: $dev_no_creds_status (should be zero)"

echo "Testing with test credentials..."
cp .env.test .env
npm run validate-env:prod
prod_with_creds_status=$?
echo "Production mode with credentials exit code: $prod_with_creds_status (should be zero)"

npm run validate-env:dev
dev_with_creds_status=$?
echo "Development mode with credentials exit code: $dev_with_creds_status (should be zero)"

echo ""
echo "Summary:"
echo "- Production mode without credentials: $([ $prod_no_creds_status -ne 0 ] && echo "PASSED (failed as expected)" || echo "FAILED (did not fail)")"
echo "- Development mode without credentials: $([ $dev_no_creds_status -eq 0 ] && echo "PASSED (continued as expected)" || echo "FAILED (did not continue)")"
echo "- Production mode with credentials: $([ $prod_with_creds_status -eq 0 ] && echo "PASSED (continued as expected)" || echo "FAILED (did not continue)")"
echo "- Development mode with credentials: $([ $dev_with_creds_status -eq 0 ] && echo "PASSED (continued as expected)" || echo "FAILED (did not continue)")"
