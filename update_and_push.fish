#!/usr/bin/env fish

function update_and_push
    set commit_message $argv[1]

    if test -z "$commit_message"
        echo "Please provide a commit message."
        exit 1
    end

    # Get the directory of this script
    set script_dir (dirname (status --current-filename))

    # Change to the script directory
    cd $script_dir

    git add .
    git commit -m "$commit_message"
    git push

    # Assuming dirdiag is in the script directory
    go install ./...
    echo "Update and push completed."

    # Change back to the original directory
    cd -
end

update_and_push $argv