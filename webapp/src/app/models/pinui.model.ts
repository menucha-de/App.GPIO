
import { Field } from './field.model';

export class PinUI extends Field {
    constructor(source?: Field) {
        super(source);
        if (source) {
            Object.assign(this, source);
            if (this.properties.get('initialState') === undefined || this.properties.get('initialState') === 'UNKNOWN') {
                this.initialState = null;
            }
        }
    }

    set status(val: boolean) {
        if (val) {
            this.value = 'HIGH';
        } else {
            this.value = 'LOW';
        }
    }

    get status() {
        if (this.value === 'HIGH') {
            return true;
        } else {
            return false;
        }
    }
    get direction() {
        return this.properties.get('direction');
    }
    set direction(value: string) {
        this.properties.set('direction', value);
    }

    get initialState() {
        return this.properties.get('initialState');
    }
    set initialState(value: string) {
        this.properties.set('initialState', value);
    }
}
