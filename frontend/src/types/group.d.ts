
interface IGroup {
    id:         number;
    groupName:  string;
    usersIds:   number[];
    users?:     IUser[];
    groupRule:  IRule[];
}