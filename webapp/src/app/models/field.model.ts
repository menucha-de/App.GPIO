export class Field {
    id?: string;
    name?: string;
    label?: string;
    value?: string;
    properties?: Map<string, string>;
    constructor(source?: Field) {
        if (source) {
            Object.assign(this, source);
            this.properties = new Map<string, string>();
            if (source['properties']) {
                this.setProperties(source['properties']);
            }
        }
    }
    setProperties(value: object) {

        if (!this.properties) {
            this.properties = new Map<string, string>();
        } else {
            this.properties.clear();
        }
        for (const [pkey, pvalue] of Object.entries(value)) {
            this.properties.set(pkey, pvalue as string);
        }
        if (this.properties.get('initialState') === undefined || this.properties.get('initialState') === 'UNKNOWN') {
            this.properties.set('initialState', null);
        }
    }
    equals(v: Field): boolean {
        if (this.id != v.id) {
            return false;
        }
        if (this.value != v.value) {
            return false;
        }
        if (this.name != v.name) {
            return false;
        }
        if (this.label != v.label) {
            return false;
        }
        if (this.properties.size !== v.properties.size) {
            return false;
        }
        let eq = true;
        this.properties.forEach((val, k) => {

            if (!v.properties.has(k)) {
                eq = false;
            }
            if (val != v.properties.get(k)) {
                eq = false;
            }
        });
        return eq;
    }
    toObject() {
        const props = Array.from(this.properties.entries()).reduce((main, [key, val]) => ({ ...main, [key]: val }), {});
        return {
            id: this.id, name: this.name,
            label: this.label, properties: props
        };
    }
}
