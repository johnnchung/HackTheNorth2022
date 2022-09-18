import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import "./Modal.css";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCircleInfo,
  faClose,
  faVideo,
} from "@fortawesome/free-solid-svg-icons";
import LoadingIcons from "react-loading-icons";

function Modal({ phrase, setModal, ...props }) {
  const [videoLink, setVideoLink] = useState("");
  const axios = require("axios").default;

  const closeModal = (e) => {
    setModal(false);
  };
  useEffect(() => {
    axios
      .post("http://localhost:8080/api/v1/process", {
        text: phrase,
      })
      .then(function ({ data }) {
        setVideoLink(data["video_url"]);
      })
      .catch((err) => {
        toast.error(err);
        return;
      });
  }, []);

  return (
    <div className="modal-container" onClick={closeModal}>
      <div className="modal-box" onClick={(e) => e.stopPropagation()}>
        <div className="modal-bar">
          <FontAwesomeIcon icon={videoLink === "" ? faCircleInfo : faVideo} />
        </div>
        <div className="modal-content">
          {videoLink !== "" ? (
            <SuccessMessage videoLink={videoLink} />
          ) : (
            <LoadingMessage />
          )}
          <p className="info-box">
            To exit and <span>discard current data</span>, please press anywhere
            outside the box.
            <br />
            <span>Remember!</span> Once closed, links cannot be retrieved again
            and will have to be generated from the start!
          </p>
        </div>
      </div>
    </div>
  );
}

function SuccessMessage({ videoLink }) {
  return (
    <div className="success-container">
      <h2>
        ✂ Your <span>SNAPPIT</span> is ready!{" "}
      </h2>
      <iframe
        src={videoLink}
        frameborder="0"
        allow="autoplay; encrypted-media"
        allowfullscreen
        title="video"
        width="565"
        height="280"
      />
      <a href={videoLink} className="dwn-btn">
        DOWNLOAD
      </a>
    </div>
  );
}

function LoadingMessage() {
  return (
    <div className="loading-container">
      <h1>Hello!</h1>
      <p className="sub-head">
        We are currently processing your data and will get back with your{" "}
        <span>SNIPPIT</span> soon...
      </p>
      <LoadingIcons.Circles stroke="#2980b9" />
      <p className="wait-time">
        <span>Waiting time:</span> 2-5 mins
      </p>
    </div>
  );
}

export default Modal;
