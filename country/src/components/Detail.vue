<template>
    <div>
        <!--    头部-->
        <div class="header">
            <div class="header-content">
                <div class="header-img-big">
                    <img class="header-img" src="../assets/header-img.png"/>
                </div>
                <div class="header-house-big">
                    <img class="header-house" src="../assets/header-house.png"/>
                </div>
            </div>
            <div class="header-text">
                <div class="header-text-first">出入境管理系统</div>
                <div class="header-text-second">Immigration management system</div>
            </div>
        </div>
        <!--    导航栏-->
        <div class="nav">
            <div class="nav-content">
                <div @click.stop="goHome()" class="nav-list">
                    <div @click="navline=0" :class="navline==0 ? 'nav-text':'nav-text-second'">申请</div>
                    <div v-show="navline==0" class="nav-line"></div>
                </div>
                <div @click.stop="goDetail()" class="nav-list">
                    <div @click="navline=2" :class="navline==2 ? 'nav-text':'nav-text-second'">详情</div>
                    <div v-show="navline==2" class="nav-line"></div>
                </div>
                <div class="nav-list">
                    <div @click="navline=1" :class="navline==1 ? 'nav-text':'nav-text-second'">查询</div>
                    <div v-show="navline==1" class="nav-line"></div>
                </div>

                <div class="nav-list">
                    <div @click="goScan()" :class="navline==3 ? 'nav-text':'nav-text-second'">区块链浏览器</div>
                    <div v-show="navline==3" class="nav-line"></div>
                </div>

            </div>
        </div>
        <!--    内容部分-->
        <div class="section">
            <div class="section-content">
                <!--        内容-->
                <div class="section-list">
                    <!--            姓名-->
                    <div class="section-tab">
                        <div class="section-text">&nbsp;Txhash：</div>
                        <div style="width: 360px;word-wrap:break-word;" class="section-input">{{BHash.hash}}</div>
                    </div>
                    <div class="section-tab">
                        <div class="section-text"><span>高</span> 度：</div>
                        <div style="width: 360px;word-wrap:break-word;" class="section-input">{{BHash.blockNumber}}</div>
                    </div>
                    <div class="section-tab">
                        <div class="section-text">高度hash：</div>
                        <div style="width: 360px;word-wrap:break-word;" class="section-input">{{BHash.blockHash}}</div>
                    </div>
                    <div class="section-tab">
                        <div class="section-text">对链txhash：</div>
                        <div style="width: 360px;word-wrap:break-word;" class="section-input">{{receipt.transactionHash}}</div>
                    </div>

                    <div class="section-tab">
                        <div class="section-text">对链高度：</div>
                        <div style="width: 360px;word-wrap:break-word;" class="section-input">{{receipt.blockNumber}}</div>
                    </div>
                    <div class="section-tab">
                        <div class="section-text">对链高度hash：</div>
                        <div style="width: 360px;word-wrap:break-word;" class="section-input">{{receipt.blockHash}}</div>
                    </div>
                </div>
                <!--        身份证-->
                <div style="margin-left: 100px;padding-top: 50px">
                    <div class="section-identity">
                        <!--        内容-->
                        <div class="section-list" style="margin-left: 150px">
                            <!--            姓名-->
                            <div class="section-tab">
                                <div class="section-text-a"><span>姓</span>名：</div>
                                <div class="section-input">{{userInfo.name}}</div>
                            </div>
                            <div class="section-tab">
                                <div class="section-text-a"><span>民</span> 族：</div>
                                <div class="section-input">{{userInfo.nation}}</div>
                            </div>
                            <div class="section-tab">
                                <div class="section-text-a">身份证号：</div>
                                <div style="width: 200px;word-wrap:break-word;" class="section-input">{{userInfo.id}}</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="loading" v-show="loading==true">
            <img src="../assets/loading.gif"/>
        </div>
    </div>
</template>

<script>
    import Web3 from "web3";
    import BABI from '../libs/bABI'
    export default {
        name: "Detail",
        data() {
            return {
                navline: 2,
                loading: true,
                timer: '',
                privateKey: '',
                account: {
                    address: ''
                },
                contract: '0x5221538292372CA706a62EB0D4890d2C85543bE9',
                providerA: 'http://39.98.214.163:8547',
                providerB: 'http://47.92.225.192:8547',
                userInfo: {
                    name: '憨憨',
                    nation: '汉族',
                    number: '1212312312312321',
                    address: '地址地址地址地址地址地址地址地址地址',
                },
                receipt:{

                },
                BHash:{
                    blockHash:'',
                    blockNumber:'',
                    hash: ''
                }
            }
        },
        mounted() {
            this.init()
        },
        methods: {
            init() {
                const web3A = new Web3(this.providerA)
                this.web3A = web3A

                const web3B = new Web3(this.providerB)
                this.web3B = web3B

                let masterPrivateKey = localStorage.getItem('masterPrivateKey')
                this.transactionHash = localStorage.getItem('transactionHash')
                if (masterPrivateKey) {
                    this.privateKey = masterPrivateKey
                    this.initAccount()
                }
            },
            async initAccount() {
                try {
                    const accountA = this.web3A.eth.accounts.privateKeyToAccount(this.privateKey)
                    this.account = accountA
                    this.web3A.eth.accounts.wallet.add(accountA)

                    const accountB = this.web3A.eth.accounts.privateKeyToAccount(this.privateKey)
                    this.account = accountB
                    this.web3B.eth.accounts.wallet.add(accountB)

                    localStorage.setItem('masterPrivateKey', this.privateKey)
                } catch (error) {
                    this.$Message.error('无效私钥，请重新输入！')
                }
                this.bindData();
            },
            async bindData(){
                console.log(`transactionHash`,this.transactionHash)
                let receipt = await this.web3A.eth.getTransactionReceipt(this.transactionHash);
                this.receipt = receipt
                console.log(`receipt`,receipt)
                this.timer = setInterval(this.getHashData, 1000);

            },
            async getHashData(){
                //获取数据
                let postObj = {
                    "method":"mos_getOtherTransactionByHash",
                    "params":[this.transactionHash],
                    "id":1
                }
                let result = await this.$http.getDetail(postObj);
                if(result && result.result){
                    clearInterval(this.timer);
                    this.BHash = result.result
                    this.BHash.blockNumber = Number(this.BHash.blockNumber).toString(10)
                    //查询用户信息
                    const contract = new this.web3B.eth.Contract(BABI)
                    contract.options.address = this.contract
                    console.log(`address`,this.account.address)
                    let user = await contract.methods.userInfo(this.account.address).call()
                    console.log(`user`,user)
                    this.userInfo = user;
                    this.loading = false;

                }
                console.log(`result`,result)
            },
            goHome() {
                this.navline = 0;
                this.$router.push({path: '/home'})
            },
            goDetail() {
                this.navline = 2;
                this.$router.push({path: '/detail'})
            },
            goScan(){
                window.open('http://39.98.50.87:8001/')
            }
        }
    }
</script>

<style scoped lang="less">
    .hello {
        position: relative;
    }

    .header {
        width: 100%;
        height: 100%;
        position: relative;
        z-index: 50;
        height: 141px;
        background: linear-gradient(354deg, #F7F7DF 0%, #0F8ECE 100%);

        .header-content {
            display: flex;
            flex-direction: row;
            flex-wrap: wrap;

            .header-img-big {
                width: 60%;
                height: 141px;

                .header-img {
                    margin-left: 20%;
                    width: 100%;
                    height: 100%;
                }
            }

            .header-house-big {
                width: 30%;
                height: 141px;

                .header-house {
                    margin-left: 10%;
                    width: 100%;
                    height: 100%;
                }
            }
        }

        .header-text {
            display: flex;
            flex-direction: column;
            align-items: flex-start;
            position: absolute;
            top: 30%;
            left: 20%;

            .header-text-first {
                font-size: 34px;
                font-family: PingFangSC-Medium, PingFang SC;
                font-weight: 500;
                color: #FFFFFF;
                line-height: 48px;
                letter-spacing: 17px;
            }

            .header-text-second {
                font-size: 19px;
                font-family: NexaBold;
                color: #FFFFFF;
                line-height: 23px;

            }
        }
    }

    .nav {
        z-index: 1000;
        background: linear-gradient(180deg, #1691F9 0%, #0070DC 100%);

        .nav-content {
            padding-left: 20%;
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: flex-start;
            padding-top: 15px;
            padding-bottom: 15px;

            .nav-list {
                display: flex;
                flex-direction: column;
                align-items: center;

                .nav-text-second {
                    cursor: pointer;
                    padding-bottom: 3px;
                    width: 100%;
                    padding-left: 50px;
                    font-size: 18px;
                    font-family: PingFangSC-Regular, PingFang SC;
                    font-weight: 400;
                    color: #FFFFFF;
                    line-height: 25px;
                }

                .nav-text {
                    cursor: pointer;
                    padding-bottom: 3px;
                    width: 100%;
                    padding-left: 50px;
                    display: flex;
                    font-size: 18px;
                    font-family: PingFangSC-Semibold, PingFang SC;
                    font-weight: 600;
                    color: #FFFFFF;
                    line-height: 25px;
                }

                .nav-line {
                    //margin-left: 20px;
                    width: 36px;
                    height: 2px;
                    background: #FFFFFF;

                }
            }
        }
    }

    .section {
        background-color: #EFF8FF;

        .section-content {
            display: flex;
            flex-direction: row;
            align-items: center;
            padding-left: 20%;
            padding-right: 20%;
            background-color: #fff;

            .section-list {
                margin-top: 100px;
                //padding-top: 45px;
                //margin-left: 150px;
                padding-left: 25px;

                .section-btn {
                    cursor: pointer;
                    width: 115px;
                    height: 37px;
                    background: linear-gradient(360deg, #0679E4 0%, #16B8FE 100%);
                    border-radius: 23px;
                    color: white;
                    border: 0;
                    outline: none;
                }

                .section-tab {
                    padding-top: 40px;
                    width: 100%;
                    display: flex;
                    flex-direction: row;
                    align-items: center;

                    .section-text {
                        width: 130px;
                        font-size: 16px;
                        font-family: PingFangSC-Semibold, PingFang SC;
                        font-weight: 600;
                        color: #1383FF;
                        line-height: 22px;

                        span {
                            padding-right: 25px
                        }
                    }

                    .section-input {
                        width: 75%;

                        .section-input-text {
                            padding-left: 10px;
                            width: 100%;
                            height: 41px;
                            background: #FFFFFF;
                            box-shadow: 2px 2px 3px 0px rgba(0, 0, 0, 0.11);
                            border-radius: 4px;
                            border: 1px solid #A7D5FF;
                            outline: none;
                        }

                        .section-input-address {
                            padding-top: 30px;
                            padding-left: 10px;
                            width: 100%;
                            height: 66px;
                            background: #FFFFFF;
                            box-shadow: 2px 2px 3px 0px rgba(0, 0, 0, 0.11);
                            border-radius: 4px;
                            border: 1px solid #A7D5FF;

                        }
                    }

                }
            }

            .section-text-a {
                width: 100px;
                color: #DC937B;
            }

            .section-identity {
                background-image: url('../assets/indentity.png');
                background-size: 100% 100%;
                width: 482px;
                height: 294px;

                img {
                    width: 440px;
                    height: 264px;
                }
            }
        }
    }
    .loading {
        position: absolute;
        top: 0;
        right: 0;
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 100%;
        height: 100%;
        background: #000;
        opacity: 0.1;

        img {
            margin-top: 300px;
        }
    }
</style>
