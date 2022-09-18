import { AwesomeButton } from "react-awesome-button";
import "react-awesome-button/dist/styles.css";
import confetti from "../assets/confetti.gif";
import "./Text.css"

function Button() {
  return (
    <AwesomeButton type="secondary" className="generate-button">
      <img src={confetti} />
      Generate Video
      <img src={confetti} />
    </AwesomeButton>
  );
}

export function TextArea() {
  return (
    <div className="container" style={{ cursor: `url(assets/ew.png), auto` }}>
      <h1>enter your phrase here</h1>
      <textarea></textarea>
      <Button />
    </div>
  );
}
