import { Field } from './field.model';

export class Device {
    constructor(source?: Device) {
        if (source) {
            Object.assign(this, source);
            this.properties = source['properties'];
            this.fields = new Map<string, Field>();
            this.setFields(source['fields']);


        } else {
            this.properties = new Map<string, string>();
            this.fields = new Map<string, Field>();
            this.id = null;
        }
    }
    id?: string;
    name?: string;
    label?: string;
    usable?: boolean;
    customized?: boolean;
    properties?: Map<string, string>;
    fields?: Map<string, Field>;
    setProperties(value: object) {
        if (value == null || value === undefined) {
            return;
        }
        if (!this.properties) {
            this.properties = new Map<string, string>();
        } else {
            this.properties.clear();
        }
        for (const [pkey, pvalue] of Object.entries(value)) {
            this.properties.set(pkey, pvalue as string);
        }
    }
    setFields(value: object) {
        if (value == null || value === undefined) {
            return;
        }
        if (!this.fields) {
            this.fields = new Map<string, Field>();
        } else {
            this.fields.clear();
        }
        for (const [pkey, pvalue] of Object.entries(value)) {
            const aa = pvalue as Field;
            const ff = new Field(pvalue);
            ff.properties = new Map(Object.entries(aa.properties));
            this.fields.set(pkey, ff);
        }

    }
    toObject() {
        const props = Array.from(this.properties.entries()).reduce((main, [key, val]) => ({ ...main, [key]: val }), {});
        const xx = {};


        for (const [key, field] of this.fields) {
            const fprops = Array.from(field.properties.entries()).reduce((main, [kkey, val]) => ({ ...main, [kkey]: val }), {});
            xx[key] = { id: field.id, name: field.name, label: field.label, properties: fprops };
        }
        return {
            id: this.id, name: this.name,
            label: this.label, properties: props,
            fields: xx
        };
    }
}
