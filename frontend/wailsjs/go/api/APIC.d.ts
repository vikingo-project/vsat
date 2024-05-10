// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {models} from '../models';
import {api} from '../models';
import {manager} from '../models';

export function About():Promise<{[key: string]: any}>;

export function CreateService(arg1:models.WebService):Promise<string>;

export function FileTypes():Promise<api.RecordsContainer>;

export function Files(arg1:string):Promise<api.RecordsContainer>;

export function GetManager():Promise<manager.Manager>;

export function Modules():Promise<api.RecordsContainer>;

export function Networks():Promise<api.RecordsContainer>;

export function RemoveService(arg1:models.ServiceHash):Promise<void>;

export function Services(arg1:string):Promise<api.RecordsContainer>;

export function SessionEvents(arg1:string):Promise<api.RecordsContainer>;

export function Sessions(arg1:string):Promise<api.RecordsContainer>;

export function ToggleService(arg1:models.ChangeServiceState):Promise<void>;

export function Tunnels(arg1:string):Promise<api.RecordsContainer>;

export function UpdateService(arg1:models.WebService):Promise<string>;
