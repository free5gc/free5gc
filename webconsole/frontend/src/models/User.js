import Serializable from "./Serializable";

export default class User extends Serializable{
  username = "";
  name = "";
  imageUrl = "";

  constructor(username, name) {
    super();
    this.username = username;
    this.name = name;
    this.imageUrl = 'https://cdn1.iconfinder.com/data/icons/evil-icons-user-interface/64/avatar-256.png';
  }
}
