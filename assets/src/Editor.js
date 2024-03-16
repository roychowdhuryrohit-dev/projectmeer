import './Editor.css'
import { $createTextNode, $getRoot, $createParagraphNode } from 'lexical';
import { useEffect } from 'react';
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
      console.log(getCanRefresh())
      if (getCanRefresh()) {
        fetch(host + "/web/getText")
          .then(res => res.json())
          .then((res) => {
            if (setPreviousText(res)) {
              editor.update(() => {
                const root = $getRoot();
                root.clear();
                const p = $createParagraphNode();
                p.append($createTextNode(res));
                root.append(p);

              });
            }
          })
          .catch((error) => {
            console.log(error)
          });
      }
    }, 2000);

    return () => clearInterval(interval);
  }, []);
}

export default function Editor() {
  let previousText = null;
  let canRefresh = true;

  const setPreviousText = (t) => {
    if (previousText === t) {
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
      console.log('ERROR:', e)
    }
  };

  const updateEditorState = (editorState) => {

    editorState.read(() => {
      setCanRefresh(false)
      if (previousText === null || previousText === undefined) {
        setPreviousText($getRoot().getTextContent());
        setCanRefresh(true)
        return;
      }
      if (previousText === $getRoot().getTextContent()) {
        setCanRefresh(true)
        return;
      }
      const newText = $getRoot().getTextContent();
      const diff = JsDiff.diffChars(previousText, newText);


      for (let d of diff) {
        console.log(d);

        if (d.added) {
          const startIndex = getChangeStartIndex(diff, d);
          console.log("Start index in old text:", startIndex);
          insertTextApi(startIndex, d.value)

        } else if (d.removed) {
          const startIndex = getChangeStartIndex(diff, d);
          console.log("Start index in old text:", startIndex);
          deleteTextApi(startIndex, 1)
        }
      }
      setPreviousText(newText);
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
}

function Placeholder() {
  return <div className="editor-placeholder">Enter some text...</div>;
}