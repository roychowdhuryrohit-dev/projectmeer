const JsDiff = require('diff');

const oldText = "The quick brown fox.";
const newText = "The quick brown dog jumps over the lazy cat.";

const diff = JsDiff.diffChars(oldText, newText);

// Function to calculate the start index of a change object
export function getChangeStartIndex(change) {
    let index = 0;

    for (let i = 0; i < diff.length; i++) {
        if (diff[i] === change) {
        return index;
        }   

        if (!diff[i].removed) {
            index += diff[i].value.length;
        }
    }

    return -1; // Change not found
}

// // Example usage:
//     console.log(oldText);
//     console.log(newText);

//     for (const d of diff) {
//         if (d.added!== undefined || d.removed!== undefined) {
//         console.log(d);
//         const startIndex = getChangeStartIndex(d);
//         console.log("Start index in old text:", startIndex);
// //After this, get new text from GET request and update state
//     }
// }

