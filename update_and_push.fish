#!/usr/bin/env fish

function update_and_push
    set commit_message $argv[1]

    if test -z "$commit_message"
        echo "Please provide a commit message."
        exit 1
    end

    git add .
    git commit -m "$commit_message"
    git push

    # Assuming dirdiag is in the current directory
    go install ./...
    echo "Update and push completed."
end

update_and_push $argv
