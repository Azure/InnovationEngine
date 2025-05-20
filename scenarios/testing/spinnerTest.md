# Testing Spinner with Elapsed Time

This test demonstrates the spinner with elapsed time display feature. The command below will run for at least 5 seconds, allowing you to observe the spinner animation with the elapsed time counter.

## Execute a long-running command

The spinner animation will display an elapsed time counter for operations that take more than a moment to complete.

```bash
echo "Starting a long-running operation..."
sleep 5
echo "Operation completed successfully!"
```

<!-- expected_similarity=1.0 -->

```text
Starting a long-running operation...
Operation completed successfully!
```

## Execute another long-running command with progress messages

This demonstrates a longer-running command with multiple progress updates. The spinner with elapsed time will be displayed between each update.

```bash
echo "Starting a complex operation with progress updates..."
for i in {1..5}; do
    echo "Processing step $i of 5..."
    sleep 1
done
echo "All steps completed!"
```

<!-- expected_similarity=1.0 -->

```text
Starting a complex operation with progress updates...
Processing step 1 of 5...
Processing step 2 of 5...
Processing step 3 of 5...
Processing step 4 of 5...
Processing step 5 of 5...
All steps completed!
```