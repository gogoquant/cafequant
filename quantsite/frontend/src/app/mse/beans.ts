export class AppInfo {
    id: number;
    name: string;
    appId: string;
    orgId: string;
    orgName: string;
    ownerName: string;
    ownerEmail: string;
    isDelete: boolean;
    dataChangeCreatedBy: string;
    dataChangeCreatedTime: string;
    dataChangeLastModifiedBy: string;
    dataChangeLastModifiedTime: string;
}

export class Organization {
    orgId: string;
    orgName: string;
}
