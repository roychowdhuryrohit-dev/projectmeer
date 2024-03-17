import './Editor.css'
import { $createTextNode, $getRoot, $createParagraphNode } from 'lexical';
import { useEffect, memo } from 'react';
import { LexicalComposer } from '@lexical/react/LexicalComposer';
import { PlainTextPlugin } from '@lexical/react/LexicalPlainTextPlugin';
import { ContentEditable } from '@lexical/react/LexicalContentEditable';
import { OnChangePlugin } from '@lexical/react/LexicalOnChangePlugin';
import { useLexicalComposerContext } from '@lexical/react/LexicalComposerContext';
const JsDiff = require('diff');

const host = "http://" + window.location.host


const insertTextApi = async (idx, vl) => {
  const requestOptions = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ index: idx, value: vl })
  };
  const response = await fetch(host + "/web/insertText", requestOptions);
  return response.ok
};

const deleteTextApi = async (idx, ct) => {
  const requestOptions = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ index: idx, count: ct })
  };
  const response = await fetch(host + "/web/deleteText", requestOptions);
  return response.ok
};

const getChangeStartIndex = (diff, change) => {
  let index = 0;

  for (let i = 0; i < diff.length; i++) {
    if (diff[i] === change) {
      return index;
    }

    if (!diff[i].removed) {
      index += diff[i].value.length;
    }
  }

  return -1;
}

const RefreshTextPlugin = ({ setPreviousText, getCanRefresh }) => {

  const [editor] = useLexicalComposerContext();
  useEffect(() => {
    const interval = setInterval(() => {
      if (getCanRefresh()) {
        fetch(host + "/web/getText")
          .then(res => res.json())
          .then((res) => {
            editor.update(() => {
              const root = $getRoot();
              if (setPreviousText(res)) {
                // console.log("updating editor", root.getTextContent(), res)
                const p = $createParagraphNode();
                p.append($createTextNode(res));
                root.clear();
                root.append(p);
                // console.log("new value", $getRoot().getTextContent());

              }
            });
          })
          .catch((error) => {
            window.location.reload()
          });
      }
    }, 1000);

    return () => clearInterval(interval);
  }, []);
}

export default memo(function Editor() {
  let previousText = null;
  let canRefresh = true;

  const setPreviousText = (t) => {
    if (previousText == t) {
      return false
    }
    previousText = t
    return true
  }
  const setCanRefresh = (val) => { canRefresh = val }
  const getCanRefresh = () => canRefresh;


  const editorConfig = {
    namespace: 'Meer',
    onError: (e) => {
      console.log(e)
    }
  };

  const updateEditorState = (editorState) => {

    editorState.read(async () => {
      setCanRefresh(false)
      let newText = $getRoot().getTextContent();
      if (previousText === null || previousText === undefined) {
        previousText = newText.trim()
        // console.log("prevtext first value", previousText)
        setCanRefresh(true)
        return;
      }
      // console.log(previousText, " - ", newText)
      if (previousText === newText) {
        setCanRefresh(true)
        return;
      }
      const diff = JsDiff.diffChars(previousText, newText);

      for (let d of diff) {
        // console.log(d);
        let dValue = d.value
        if (d.added) {
          if (d.value.length != 1) {
            dValue = d.value.slice(0, -2)
            newText = previousText + dValue
          }
          const startIndex = getChangeStartIndex(diff, d);
          // console.log("Insertion - Start index in old text:", startIndex);
          await insertTextApi(startIndex, dValue)

        } else if (d.removed) {
          const startIndex = getChangeStartIndex(diff, d);
          // console.log("Deletion - Start index in old text:", startIndex);
          await deleteTextApi(startIndex, d.count)
        }
      }
      // console.log("refresh true")
      previousText = newText
      setCanRefresh(true)
    });
  };


  return (
    <LexicalComposer initialConfig={editorConfig}>
      <div className="editor-container">
        <PlainTextPlugin
          contentEditable={<ContentEditable className="editor-input" />}
          placeholder={<Placeholder />}
        />
        <OnChangePlugin onChange={updateEditorState} />
        <RefreshTextPlugin setPreviousText={setPreviousText} getCanRefresh={getCanRefresh} />
      </div>
    </LexicalComposer>
  );
});

function Placeholder() {
  return <div className="editor-placeholder">Enter some text...</div>;
}