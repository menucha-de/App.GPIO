export class CycleSpec {
    id?: string;
    applicationId?: string;
    name?: string;
    enabled?: boolean;
    duration = -1;
    repeatPeriod = -1;
    interval = -1;
    whenDataAvailable = false;
    whenDataAvailableDelay = 0;
    reportIfEmpty = false;
    fieldSubscriptions: Map<string, Set<string>>;
    constructor(source?: CycleSpec) {
        if (source) {
            Object.assign(this, source);
            this.fieldSubscriptions = new Map();
            Object.keys(source.fieldSubscriptions).forEach(k => { this.fieldSubscriptions.set(k, source.fieldSubscriptions[k]); });


        } else {
            this.duration = 1000;
            this.repeatPeriod = 1000;
            this.id = null;
            this.applicationId = null;
            this.enabled = false;
            this.fieldSubscriptions = new Map();
        }
    }
    toObject() {
        const flds = Array.from(this.fieldSubscriptions.entries()).reduce((main, [key, val]) => ({ ...main, [key]: Array.from(val )}), {});
        return {
            id: this.id, name: this.name,
            applicationId: this.applicationId, duration: this.duration,
            repeatPeriod: this.repeatPeriod,
            whenDataAvailable: this.whenDataAvailable,
            whenDataAvailableDelay: this.whenDataAvailableDelay,
            reportIfEmpty: this.reportIfEmpty, interval: this.interval,
            enabled: this.enabled,
            fieldSubscriptions: flds,
        };
    }
}
