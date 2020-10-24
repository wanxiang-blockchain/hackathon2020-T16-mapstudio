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
                <div class="header-text-first">身份链管理系统</div>
                <div class="header-text-second">Identity chain management system</div>
            </div>
        </div>
        <!--    导航栏-->
        <div class="nav">
            <div class="nav-content">
                <div @click.stop="goHome()" class="nav-list">
                    <div @click="navline=0" :class="navline==0 ? 'nav-text':'nav-text-second'">登记</div>
                    <div v-show="navline==0" class="nav-line"></div>
                </div>
                <div @click.stop="goDetail()" class="nav-list">
                    <div @click="navline=2" :class="navline==2 ? 'nav-text':'nav-text-second'">详情</div>
                    <div v-show="navline==2" class="nav-line"></div>
                </div>
                <div  class="nav-list">
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
                        <div class="section-text" style="padding-left: 5px">&nbsp;Txhash：</div>
                        <div style="width: 360px;word-wrap:break-word;" class="section-input">{{transactionHash}}</div>
                    </div>
                    <div class="section-tab">
                        <div class="section-text"><span>高</span> 度：</div>
                        <div class="section-input">{{receipt.blockNumber}}</div>
                    </div>
                    <div class="section-tab">
                        <div class="section-text"><span>状</span> 态：</div>
                        <div class="section-input">{{receipt.status==true?'成功':'失败'}}</div>
                    </div>
                    <div class="section-tab">
                        <div class="section-text">高度hash：</div>
                        <div style="width: 360px;word-wrap:break-word;" class="section-input">{{receipt.blockHash}}</div>
                    </div>
                </div>
                <!--        身份证-->
                <div style="margin-left: 100px;padding-top: 50px">
                    <div class="section-identity">
                        <!--        内容-->
                        <div class="section-list">
                            <!--            姓名-->
                            <div class="section-tab">
                                <div class="section-text" style="padding-left: 5px"><span>姓</span>名：</div>
                                <div class="section-input">{{userInfo.name}}</div>
                            </div>
                            <div class="section-tab">
                                <div class="section-text"><span>民</span> 族：</div>
                                <div class="section-input">{{userInfo.nation}}</div>
                            </div>
                            <div class="section-tab">
                                <div class="section-text">身份证号：</div>
                                <div style="width: 150px;word-wrap:break-word;" class="section-input">{{userInfo.id}}</div>
                            </div>
                            <div class="section-tab" style="padding-top: 20px">
                                <div class="section-text"><span>地</span>址：</div>
                                <div style="width: 300px;word-wrap:break-word;" class="section-input">{{userInfo.userAddress}}</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import Web3 from "web3";
    import AABI from '../libs/aABI'

    export default {
        name: "Detail",
        data() {
            return {
                navline: 2,
                privateKey: '',
                account: {
                    address: ''
                },
                contract: '0x5221538292372CA706a62EB0D4890d2C85543bE9',
                provider: 'http://39.98.214.163:8547',
                transactionHash:'',
                userInfo: {
                    name: '憨憨',
                    nation: '汉族',
                    id: '1212312312312321',
                    userAddress: '地址地址地址地址地址地址地址地址地址',
                },
                receipt:{}
            }
        },
        mounted() {
            this.init()
        },
        methods: {
            init() {
                const web3 = new Web3(this.provider)
                this.web3 = web3
                let masterPrivateKey = localStorage.getItem('masterPrivateKey')
                this.transactionHash = localStorage.getItem('transactionHash')
                if (masterPrivateKey) {
                    this.privateKey = masterPrivateKey
                    this.initAccount()
                }

            },
            async initAccount() {
                try {
                    const account = this.web3.eth.accounts.privateKeyToAccount(this.privateKey)
                    this.account = account
                    this.web3.eth.accounts.wallet.add(account)
                    localStorage.setItem('masterPrivateKey', this.privateKey)
                } catch (error) {
                    this.$Message.error('无效私钥，请重新输入！')
                }
                this.bindData();
            },
            async bindData(){
                console.log(`transactionHash`,this.transactionHash)
                let receipt = await this.web3.eth.getTransactionReceipt(this.transactionHash);
                this.receipt = receipt
                console.log(`receipt`,receipt)
                const contract = new this.web3.eth.Contract(AABI)
                contract.options.address = this.contract
                console.log(`address`,this.account.address)
                let user = await contract.methods.userInfo(this.account.address).call()
                console.log(`user`,user)
                this.userInfo = user;
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
                window.open('http://39.98.50.87:8000/')
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
                //margin-top: 10px;
                //padding-top: 45px;
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
                        width: 100px;
                        font-size: 16px;
                        font-family: PingFangSC-Semibold, PingFang SC;
                        font-weight: 600;
                        color: #1383FF;
                        line-height: 22px;
                        padding-right: 20px;

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
</style>
