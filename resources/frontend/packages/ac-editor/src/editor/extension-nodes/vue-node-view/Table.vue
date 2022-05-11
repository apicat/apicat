<template>
    <node-view-wrapper :class="className">
        <div class="scrollable">
            <table class="rme-table">
                <colgroup>
                    <col width="33.3333%" />
                    <col width="33.3333%" />
                    <col width="33.3333%" />
                </colgroup>
                <tbody>
                    <tr>
                        <th>
                            <a class="grip-row first" contenteditable="false"></a>
                            <a class="grip-column first" contenteditable="false"></a>
                            <p contenteditable="true"></p>
                        </th>
                        <th>
                            <a class="grip-column" contenteditable="false"></a>
                            <p contenteditable="true"></p>
                        </th>
                        <th>
                            <a class="grip-column" contenteditable="false"></a>
                            <p contenteditable="true"></p>
                        </th>
                    </tr>
                    <tr>
                        <td>
                            <a class="grip-row" contenteditable="false"></a>
                            <p contenteditable="true"></p>
                        </td>
                        <td><p contenteditable="true"></p></td>
                        <td><p contenteditable="true"></p></td>
                    </tr>
                    <tr>
                        <td>
                            <a class="grip-row last" contenteditable="false"></a>
                            <p contenteditable="true" @click="onInput"></p>
                        </td>
                        <td><p contenteditable="true"></p></td>
                        <td><p contenteditable="true"></p></td>
                    </tr>
                </tbody>
            </table>
        </div>
    </node-view-wrapper>
</template>

<script>
    /**
     * 表格组件
     */
    import { NodeViewWrapper } from '../../../components/NodeViewWrapper'

    export default {
        name: 'Table',
        props: ['node', 'view', 'getPos', 'isSelected', 'isReadOnly', 'updateAttributes'],
        components: {
            NodeViewWrapper,
        },
        computed: {
            attrs: function () {
                return this.node.attrs
            },
            className() {
                return ['scrollable-wrapper original-table']
            },
        },
        methods: {
            onInput() {
                // console.log(123);
            },
        },
    }
</script>
