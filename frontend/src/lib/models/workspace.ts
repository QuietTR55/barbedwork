import type { Channel } from './channel';
import type { User } from './user';

export interface ImagePathObject {
	String: string;
	Valid: boolean;
}

export type Workspace = {
	Id: string;
	ImagePath: ImagePathObject;
	Name: string;
	Users: User[];
	Channels: Channel[];
	Permissions: string[];
};
