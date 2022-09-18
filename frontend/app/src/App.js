import "./App.css";
import { TextArea } from "./components";
import {
  CookieMonster,
  DonaldTrump,
  BenShapiro,
  Penguin,
  Giraffe,
  Logo,
} from "./components/images";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { useState } from "react";
import Modal from "./components/Modal";

function App() {
  const [modal, setModal] = useState(false);
  const [phrase, setPhrase] = useState("");

  return (
    <div className="App">
      {modal ? <Modal phrase={phrase} setModal={setModal} /> : null}
      <CookieMonster />
      <DonaldTrump />
      <BenShapiro />
      <Penguin />
      <Giraffe />
      <Logo />
      <TextArea
        modal={modal}
        setModal={setModal}
        phrase={phrase}
        setPhrase={setPhrase}
      />
      <ToastContainer />
    </div>
  );
}

export default App;
