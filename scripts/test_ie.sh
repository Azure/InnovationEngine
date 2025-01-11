for file in docs/*; do
    echo "=== Testing '$file' ==="
    ie test "$file"
done