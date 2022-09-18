import { useState } from "react";
import { AwesomeButton } from "react-awesome-button";
import "react-awesome-button/dist/styles.css";
import { toast } from "react-toastify";
import confetti from "../assets/confetti.gif";
import "./Text.css";

export function TextArea({ modal, setModal, phrase, setPhrase, ...props }) {
  const handleSubmit = () => {
    if (phrase === "") {
      toast.error("No Phrase entered! Whatcha you got on your mind?");
      return;
    }

    setModal(true);
  };

  return (
    <div className="container" style={{ cursor: `url(assets/ew.png), auto` }}>
      <h1>enter your phrase here</h1>
      <textarea
        value={phrase}
        onChange={(e) => setPhrase(e.target.value)}
      ></textarea>

      <AwesomeButton
        type="secondary"
        className="generate-button"
        action={handleSubmit}
      >
        <img src={confetti} />
        Generate Video
        <img src={confetti} />
      </AwesomeButton>
    </div>
  );
}
